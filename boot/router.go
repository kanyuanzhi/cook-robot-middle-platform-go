package boot

import (
	"github.com/gin-gonic/gin"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/middleware"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/router/v1"
)

func Router() *gin.Engine {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)
	router.Use(middleware.Cors())
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/api/v1/controller/fetch-status"))

	apiV1 := router.Group("/api/v1")

	// 不用验证的路由
	//publicGroup := apiV1.Group("/public")
	//v1.InitPublicRouter(publicGroup)

	// 需要验证的路由
	//privateGroup := apiV1.Group("/private")
	v1.InitPrivateRouter(apiV1)

	return router
}
