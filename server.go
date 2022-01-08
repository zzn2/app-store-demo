package main

import (
	"net/http"

	"example.com/goserver/app"
	"github.com/gin-gonic/gin"
)

var store app.Store

func newApp(c *gin.Context) {
	var app app.Meta

	if err := c.ShouldBindYAML(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	store.Add(app)
	c.JSON(http.StatusCreated, gin.H{
		"message": app.Title,
	})
}

func listApps(c *gin.Context) {
	c.JSON(http.StatusOK, store.List())
}

func main() {
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/apps", newApp)
		v1.GET("/apps", listApps)
	}

	router.Run(":3001")
}
