package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/madeinatria/love-all-backend/docs"
	"github.com/madeinatria/love-all-backend/internal/handlers"
	"github.com/madeinatria/love-all-backend/internal/middleware"
	"github.com/madeinatria/love-all-backend/internal/routes"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartApp() {
	router := initializeGin()
	SetUpRoutes(router)
	router.Run(fmt.Sprintf(":%s", "8001"))
}

func SetUpRoutes(router *gin.Engine) {

	docs.SwaggerInfo.BasePath = "/api/v1"
	api := router.Group("/api/v1", middleware.Authorize())
	{
		router.POST("/login", handlers.LoginHandler)
		routes.AllRoutes(api)
		router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	}

}

func initializeGin() *gin.Engine {
	router := gin.New()
	router.Use(middleware.GetZapGinConfig())
	return router
}
