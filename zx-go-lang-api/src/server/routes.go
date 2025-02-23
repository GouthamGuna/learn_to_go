package server

import (
	"github.com/GouthamGuna/zx/src/controllers"
	"github.com/GouthamGuna/zx/src/services"
	"github.com/gin-gonic/gin"
)

// Routers sets up all the routes.
func routers(r *gin.Engine) {

	controllers.Login(r)

	apiV1 := r.Group("/api/v1")
	apiV1.Use(services.AuthenticateMiddleware)

	controllers.UtilController(apiV1)
	controllers.ZXController(apiV1)
}
