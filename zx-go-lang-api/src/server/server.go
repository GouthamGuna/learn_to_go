package server

import (
	utils "github.com/GouthamGuna/zx/src/utils"
	"github.com/gin-gonic/gin"
)

func Start() {
    // Set Gin mode to ReleaseMode
    gin.SetMode(gin.ReleaseMode)

    r := gin.Default()
    routers(r)

    // Run the server
    r.Run(utils.GetServerPort())
}
