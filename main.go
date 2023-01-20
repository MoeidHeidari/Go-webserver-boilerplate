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
// @version 1.0
// @description The BEST API you have ever seen
// @host localhost:6001
// @BasePath /api/v1
// @securityDefinitions.basic  BasicAuth

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
