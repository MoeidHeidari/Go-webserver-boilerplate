package routes

import (
	"main/api/controllers"
	"main/api/currencies"
	"main/api/kubes"

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
	kubes          kubes.KubeRequest
}

// Setup Test routes
func (s TestRoutes) Setup() {
	s.logger.Info("Setting up routes")
	api := s.handler.Gin.Group("/api") //.Use(s.authMiddleware.Handler())
	{
		api.GET("/test", s.TestController.GetTest)
		api.GET("/test/:id", s.TestController.GetOneTest)
		api.GET("/currency", s.testRequest.MakeRequest)
		api.GET("/kube_get", s.kubes.GetPodInfoRequest)
		api.POST("/test", s.TestController.CreateTest)
		api.POST("/currency", s.testRequest.MakePostRequest)
		api.POST("/test/:id", s.TestController.UpdateTest)
		api.POST("/kube_add", s.kubes.CreatePodRequest)
		api.DELETE("/test/:id", s.TestController.DeleteTest)
		api.DELETE("/kube_delete/:namespace/:pod_name", s.kubes.DeletePodRequest)
	}
}

// NewTestRoutes creates new Test controller
func NewTestRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	testRequest currencies.Request,
	TestController controllers.TestController,
	kubes kubes.KubeRequest,
) TestRoutes {
	return TestRoutes{
		handler:        handler,
		logger:         logger,
		testRequest:    testRequest,
		TestController: TestController,
		kubes:          kubes,
	}
}
