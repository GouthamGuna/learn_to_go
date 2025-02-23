package controllers

import (
	"github.com/GouthamGuna/zx/src/services"
	"github.com/gin-gonic/gin"
)

func ZXController(rg *gin.RouterGroup) {

	rg.GET("/stream/:filename", services.GetDataHandler)
	rg.GET("/filenames", services.GetFilenamesHandler)
}
