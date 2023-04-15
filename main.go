package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

var con config

type config struct {
	script_path string `yaml:"path"`
	port        string `yaml:"port"`
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
	err = yaml.Unmarshal(yamlFile, con)
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
	e.Run(":" + con.port)
}

func redeploy(c *gin.Context) {
	cmd := exec.Command(con.script_path)

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
