package main

import (
	"log"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func main() {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery(), Cors())
	e.GET("build", build)
	e.Run(":25583")
}

func build(c *gin.Context) {
	cmd := exec.Command("F:/passport-ci-test/build.bat")
	// out, err := exec.Command("F:/passport-ci-test/build.bat").Output()
	// if err != nil {
	// 	log.Println(err)
	// }
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	// 	c.JSON(200, gin.H{
	// 		"out": ,
	// 		"err": ,
	// 	})
}
