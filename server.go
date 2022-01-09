package main

import (
	"net/http"

	"example.com/goserver/app"
	"example.com/goserver/filter"
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
	q := c.Request.URL.Query()
	flt, err := filter.Create(q)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	result, err := store.List(*flt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, result)
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
