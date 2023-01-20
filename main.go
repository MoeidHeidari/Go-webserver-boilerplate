package main

import (
	"main/bootstrap"
	"main/docs"
	"main/lib"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title SkyFarm
// @description The BEST API you have ever seen
// @host localhost:3000
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	go runSwagger()
	_ = godotenv.Load()
	err := bootstrap.RootApp.Execute()
	if err != nil {
		return
	}

}

func runSwagger() {
	env := lib.NewEnv()
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.Title = "Skyfarm API"
	r.Run(":" + env.SwaggerPort)
}
