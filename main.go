package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/yaml.v3"
)

var con Config

type Config struct {
	Script_path string `yaml:"path"`
	Port        string `yaml:"port"`
	Access_key  string `yaml:"access_key"`
}

type Access_key struct {
	Key string `json:"key"`
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, PUT")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var key Access_key
		err := c.ShouldBindBodyWith(&key, binding.JSON)
		if err != nil {
			log.Println(err)
			c.JSON(401, gin.H{
				"msg": err,
			})
		}
		if key.Key != con.Access_key {
			c.Abort()
		}
		c.Next()
	}
}

func main() {
	ConfigInit(&con, "./config.yaml")
	Run()
}

func ConfigInit(c interface{}, path string) {
	path = filepath.Clean(path)
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &con)
	if err != nil {
		log.Panic(err)
	}
	log.Println("config read success")
}

func Run() {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery(), Cors(), Auth())
	e.POST("alive")
	e.PUT("redeploy", redeploy)
	e.Run(":" + con.Port)
}

func redeploy(c *gin.Context) {
	cmd := exec.Command(con.Script_path)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Println(err)
		c.JSON(200, gin.H{
			"err": err,
		})
	}

	outStr, errStr := stdout.String(), stderr.String()
	log.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	c.JSON(200, gin.H{
		"execOut": outStr,
		"execErr": errStr,
	})
}
