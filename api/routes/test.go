package routes

import (
	"main/api/controllers"
	"main/api/kubes"
	"main/api/middlewares"
	"main/api/ws"
	"main/lib"

	"github.com/gin-gonic/gin"
)

// TestRoutes struct
type TestRoutes struct {
	logger         lib.Logger
	handler        lib.RequestHandler
	TestController controllers.TestController
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
		api.GET("/kube_get/:namespace", s.kubes.GetNodeInfoRequest)
		api.GET("/helm_get", s.kubes.HGetReleaseRequest)
		api.POST("/test", s.TestController.CreateTest)
		api.POST("/test/:id", s.TestController.UpdateTest)
		api.POST("/kube/add", s.kubes.CreatePodRequest)
		api.POST("/kube/create_config_map", s.kubes.CreateOrUpdateConfigMapRequest)
		api.POST("/kube/create_secret", s.kubes.CreateOrUpdateSecretRequest)
		api.POST("/kube/create_namespace", s.kubes.CreateNamespaceRequest)
		api.POST("/kube/create_pv", s.kubes.CreatePersistentVolumeRequest)
		api.POST("/kube/create_pvc", s.kubes.CreatePersistentVolumeClaimRequest)
		api.POST("/kube/create_nodeport", s.kubes.CreateNodePortRequest)
		api.POST("/helm_create", s.kubes.HCreateReleaseRequest)
		api.POST("helm_create_repo", s.kubes.HCreateRepoRequest)
		api.POST("/kube/create_role", s.kubes.CreateRoleRequest)
		api.POST("/kube/role_bind", s.kubes.CreateRoleBindingRequest)
		api.POST("/kube/create_account", s.kubes.CreateServiceAccountRequest)
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
	TestController controllers.TestController,
	kubes kubes.KubeRequest,
) TestRoutes {
	return TestRoutes{
		handler:        handler,
		logger:         logger,
		TestController: TestController,
		kubes:          kubes,
	}
}
