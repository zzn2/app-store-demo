package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zzn2/demo/appstore/app"
	"github.com/zzn2/demo/appstore/filter"
)

var store app.Store

func newApp(c *gin.Context) {
	var app app.Meta

	if err := c.ShouldBindYAML(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	store.Add(app)
	c.JSON(http.StatusCreated, app)
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
