package routes

import (
	"main/api/controllers"
	"main/api/currencies"
	"main/api/kubes"
	"main/api/middlewares"
	"main/lib"
	"main/ws"

	"github.com/gin-gonic/gin"
)

// TestRoutes struct
type TestRoutes struct {
	logger         lib.Logger
	handler        lib.RequestHandler
	TestController controllers.TestController
	testRequest    currencies.Request
	authMiddleware middlewares.JWTAuthMiddleware
	kubes          kubes.KubeRequest
	Websocket      ws.Ws
}

// Setup Test routes
func (s TestRoutes) Setup() {
	s.logger.Info("Setting up routes")
	s.handler.Gin.GET("/get_code", s.TestController.GetCode)
	api := s.handler.Gin.Group("/api").Use(s.authMiddleware.Handler())
	{
		api.GET("/test", s.TestController.GetTest)
		api.GET("/test/:id", s.TestController.GetOneTest)
		api.GET("/currency", s.testRequest.MakeRequest)
		api.GET("/kube_get", s.kubes.GetPodInfoRequest)
		api.GET("/helm_get", s.kubes.HGetReleaseRequest)
		api.POST("/test", s.TestController.CreateTest)
		api.POST("/currency", s.testRequest.MakePostRequest)
		api.POST("/test/:id", s.TestController.UpdateTest)
		api.POST("/kube_add", s.kubes.CreatePodRequest)
		api.POST("/helm", s.kubes.HCreateReleaseRequest)
		api.POST("/helm_create_repository", s.kubes.HCreateRepositoryRequest)
		api.POST("/kube/create_config_map", s.kubes.CreateOrUpdateConfigMapRequest)
		api.POST("/kube/create_secret", s.kubes.CreateOrUpdateSecretRequest)
		api.DELETE("/test/:id", s.TestController.DeleteTest)
		api.DELETE("/kube_delete/:namespace/:pod_name", s.kubes.DeletePodRequest)

	}
	r := gin.Default()
	r.GET("/", s.Websocket.MessageHandler)
	go r.Run(":12121")

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
