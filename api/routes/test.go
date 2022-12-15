package routes

import (
	"main/api/controllers"
	"main/lib"
)

// TestRoutes struct
type TestRoutes struct {
	logger         lib.Logger
	handler        lib.RequestHandler
	TestController controllers.TestController
}

// Setup Test routes
func (s TestRoutes) Setup() {
	s.logger.Info("Setting up routes")
	api := s.handler.Gin.Group("/api")
	{
		api.GET("/test", s.TestController.GetTest)
		api.GET("/test/:id", s.TestController.GetOneTest)
		api.POST("/test", s.TestController.SaveTest)
		api.POST("/test/:id", s.TestController.UpdateTest)
		api.DELETE("/test/:id", s.TestController.DeleteTest)
	}
}

// NewTestRoutes creates new Test controller
func NewTestRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	TestController controllers.TestController,
) TestRoutes {
	return TestRoutes{
		handler:        handler,
		logger:         logger,
		TestController: TestController,
	}
}
