package main

import (
	"github.com/gin-gonic/gin"
)

func testApi(c *gin.Context) {

	c.JSON(200, gin.H{
		"data": "test api works",
	})

}
func main() {
	engine := gin.New()
	engine.GET("/test", testApi)
	err := engine.Run("localhost:4000")
	if err != nil {
		return
	}

}
