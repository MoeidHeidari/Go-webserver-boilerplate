package routes

import (
	"github.com/dipeshdulal/clean-gin/api/middlewares"
	"github.com/dipeshdulal/clean-gin/lib"
	"main/api/controllers"
)

// TestRoutes struct
type TestRoutes struct {
	logger         lib.Logger
	handler        lib.RequestHandler
	TestController controllers.TestController
	authMiddleware middlewares.JWTAuthMiddleware
}

// Setup Test routes
func (s TestRoutes) Setup() {
	s.logger.Info("Setting up routes")
	api := s.handler.Gin.Group("/api").Use(s.authMiddleware.Handler())
	{
		api.GET("/Test", s.TestController.GetTest)
		api.GET("/Test/:id", s.TestController.GetOneTest)
		api.POST("/Test", s.TestController.SaveTest)
		api.POST("/Test/:id", s.TestController.UpdateTest)
		api.DELETE("/Test/:id", s.TestController.DeleteTest)
	}
}

// NewTestRoutes creates new Test controller
func NewTestRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	TestController controllers.TestController,
	authMiddleware middlewares.JWTAuthMiddleware,
) TestRoutes {
	return TestRoutes{
		handler:        handler,
		logger:         logger,
		TestController: TestController,
		authMiddleware: authMiddleware,
	}
}
