package main

import "github.com/gin-gonic/gin"

func newApp(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func listApps(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "apps",
	})
}

func main() {
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/apps", newApp)
		v1.GET("/apps", listApps)
	}

	router.Run(":8081")
}
