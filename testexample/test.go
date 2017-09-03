package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
<<<<<<< HEAD
			"message": "pong997",
=======
			"message": "pong7",
>>>>>>> 7851e0c06507873fa1cd0e62fd34692cd8bae112
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
