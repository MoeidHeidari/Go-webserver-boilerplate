package routes

import (
	"main/api/controllers"
	"main/api/currencies"

	"main/api/middlewares"
	"main/lib"
)

// TestRoutes struct
type TestRoutes struct {
	logger         lib.Logger
	handler        lib.RequestHandler
	TestController controllers.TestController
	testRequest    currencies.Request
	authMiddleware middlewares.JWTAuthMiddleware
}

// Setup Test routes
func (s TestRoutes) Setup() {
	s.logger.Info("Setting up routes")
	api := s.handler.Gin.Group("/api").Use(s.authMiddleware.Handler())
	{
		api.GET("/test", s.TestController.GetTest)
		api.GET("/test/:id", s.TestController.GetOneTest)
		api.GET("/currency", s.testRequest.MakeRequest)
		api.POST("/test", s.TestController.CreateTest)
		api.POST("/currency", s.testRequest.MakePostRequest)
		api.POST("/test/:id", s.TestController.UpdateTest)
		api.DELETE("/test/:id", s.TestController.DeleteTest)
	}
}

// NewTestRoutes creates new Test controller
func NewTestRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	testRequest currencies.Request,
	TestController controllers.TestController,
) TestRoutes {
	return TestRoutes{
		handler:        handler,
		logger:         logger,
		testRequest:    testRequest,
		TestController: TestController,
	}
}
