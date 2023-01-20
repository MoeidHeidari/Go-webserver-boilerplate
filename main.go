package main

import (
	"main/bootstrap"
	"main/lib"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title SkyFarm
// @version 1.0
// @decription The BEST API you have ever seen

// @host localhost:8080
// @BasePath /

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
	r.Run(":" + env.SwaggerPort)
}
