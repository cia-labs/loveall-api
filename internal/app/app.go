package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madeinatria/love-all-backend/docs"
	"github.com/madeinatria/love-all-backend/internal/handlers"
	"github.com/madeinatria/love-all-backend/internal/middleware"
	"github.com/madeinatria/love-all-backend/internal/routes"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// gin-swagger middleware
	// swagger embed files
)

func StartApp() {
	router := initializeGin()
	SetUpRoutes(router)
	router.Run(fmt.Sprintf(":%s", "8081"))
}

func SetUpRoutes(router *gin.Engine) {

	docs.SwaggerInfo.BasePath = "/api/v1"
	// v1 := router.Group("/api/v1")
	// {
	// 	eg := v1.Group("/example")
	// 	{
	// 		eg.GET("/helloworld", Helloworld)
	// 	}
	// }

	api := router.Group("/api/v1", middleware.Authorize())
	{
		// api.GET("/", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"msg": ""}) })
		router.POST("/login", handlers.LoginHandler)
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		routes.AllRoutes(api)
	}

}

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}
func initializeGin() *gin.Engine {
	router := gin.New()
	router.Use(middleware.GetZapGinConfig())
	return router
}
