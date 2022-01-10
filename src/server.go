package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zzn2/demo/appstore/app"
	"github.com/zzn2/demo/appstore/filter"
)

var store app.Store

func newApp(c *gin.Context) {
	var app app.Meta

	if err := c.ShouldBindYAML(&app); err != nil {
		c.JSON(http.StatusBadRequest, responseBodyForError(err))
		return
	}
	if store.GetByTitleAndVersion(app.Title, app.Version) != nil {
		c.JSON(http.StatusBadRequest, responseBodyForErrorMessage("App '%s' with version '%s' already exists.", app.Title, app.Version))
	} else {
		store.Add(app)
		c.JSON(http.StatusCreated, app)
	}
}

func getAppByTitle(c *gin.Context) {
	title := c.Param("title")
	app := store.GetByTitle(title)
	if app != nil {
		c.JSON(http.StatusOK, app)
	} else {
		c.JSON(http.StatusNotFound, responseBodyForErrorMessage("App with title '%s' does not exist.", title))
	}
}

func listApps(c *gin.Context) {
	q := c.Request.URL.Query()
	flt, err := filter.Create(q)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseBodyForError(err))
	}

	result, err := store.List(*flt)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseBodyForError(err))
	}
	c.JSON(http.StatusOK, result)
}

// errorResponse replys a general error response
func responseBodyForError(err error) map[string]interface{} {
	return responseBodyForErrorMessage(err.Error())
}

func responseBodyForErrorMessage(format string, a ...interface{}) map[string]interface{} {
	return gin.H{"error": fmt.Sprintf(format, a...)}
}

func main() {
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/apps", newApp)
		v1.GET("/apps", listApps)
		v1.GET("/apps/:title", getAppByTitle)
	}

	router.Run(":3001")
}
