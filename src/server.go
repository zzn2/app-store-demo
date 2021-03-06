package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zzn2/demo/appstore/app"
	"github.com/zzn2/demo/appstore/filter"
	"github.com/zzn2/demo/appstore/semver"
)

var store *app.Store

func newApp(c *gin.Context) {
	var app app.Meta

	if err := c.ShouldBindYAML(&app); err != nil {
		c.JSON(http.StatusBadRequest, responseBodyForError(err))
		return
	}

	if app.Version == semver.Empty {
		err := errors.New(fmt.Sprintf("App '%s' lacks of version or the version could not be '%s'.)", app.Title, app.Version))
		c.JSON(http.StatusBadRequest, responseBodyForError(err))
		return
	}

	err := store.Add(app)
	if err != nil {
		c.JSON(http.StatusConflict, responseBodyForError(err))
	} else {
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

func getAppByTitleAndVersion(c *gin.Context) {
	title := c.Param("title")
	versionText := c.Param("version")
	version, err := semver.Parse(versionText)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseBodyForErrorMessage("Bad format of version '%s'", versionText))
	}
	app := store.GetByTitleAndVersion(title, version)
	if app != nil {
		c.JSON(http.StatusOK, app)
	} else {
		c.JSON(http.StatusNotFound, responseBodyForErrorMessage("App with title '%s' and version '%s' does not exist.", title, version))
	}
}

func listApps(c *gin.Context) {
	var app app.Meta
	q := c.Request.URL.Query()
	flt, err := filter.CreateRuleSet(q, app)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseBodyForError(err))
		return
	}

	result, err := store.List(flt)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseBodyForError(err))
		return
	}
	c.JSON(http.StatusOK, result)
}

func responseBodyForError(err error) map[string]interface{} {
	return responseBodyForErrorMessage(err.Error())
}

// responseBodyForErrorMessage is aimed to format the response body of bad requests.
// All the bad requests will go through this function so that the error responses will have the same format.
// This makes the consumer of our APIs easier to write error handling logic.
func responseBodyForErrorMessage(format string, a ...interface{}) map[string]interface{} {
	return gin.H{"error": fmt.Sprintf(format, a...)}
}

func setupServer() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/apps", newApp)
		v1.GET("/apps", listApps)
		v1.GET("/apps/:title", getAppByTitle)
		v1.GET("/apps/:title/versions/:version", getAppByTitleAndVersion)
	}

	return router
}

// This function sets up a new, empty store.
// It is supposed to be called when the app starts up.
// It could also be called from the integration test in order to get a clean store for each test scenario.
func setupStore() {
	var emptyStore app.Store
	store = &emptyStore
}

func main() {
	setupStore()
	setupServer().Run(":3001")
}
