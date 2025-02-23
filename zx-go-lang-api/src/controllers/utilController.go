package controllers

import (
	utils "github.com/GouthamGuna/zx/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UtilController(rg *gin.RouterGroup) {

	rg.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	rg.GET("/data", func(c *gin.Context) {
		// Sample JSON response
		data := gin.H{
			"message": "Hello, World!",
			"status":  "success",
		}
		c.JSON(http.StatusOK, data)
	})

	rg.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, utils.GETFavicon())
	})
}
