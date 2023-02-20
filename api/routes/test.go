package routes

import (
	"main/api/controllers"
	"main/api/kubescontrollers"
	"main/api/middlewares"
	"main/api/ws"
	"main/lib"

	"github.com/gin-gonic/gin"
)

// TestRoutes struct
type TestRoutes struct {
	logger          lib.Logger
	handler         lib.RequestHandler
	TestController  controllers.TestController
	authMiddleware  middlewares.JWTAuthMiddleware
	kubescontroller kubescontrollers.KubeController
	Websocket       ws.Ws
}

// Setup Test routes
func (s TestRoutes) Setup() {
	s.logger.Info("Setting up routes")
	s.handler.Gin.GET("/get_code", s.TestController.GetCode)
	api := s.handler.Gin.Group("/api") //.Use(s.authMiddleware.Handler())
	{
		api.GET("/workspace", s.TestController.GetWorkspaces)
		api.GET("/workspace/:id", s.TestController.GetOneWorkspace)
		api.GET("/kube_get/:namespace", s.kubescontroller.GetPodList)
		api.GET("/helm_get", s.kubescontroller.HGetReleaseRequest)
		api.GET("/workspace/trash", s.TestController.GetDeletedWorkspaces)
		api.POST("/workspace_create", s.TestController.CreateWorkspace)
		api.POST("/workspace/:id/add_edge", s.TestController.AddEdge)
		api.POST("/workspace/:id/add_node", s.TestController.AddNode)
		api.POST("/workspace/:id/update_node", s.TestController.UpdateNode)
		api.POST("/kube/add", s.kubescontroller.CreatePodRequest)
		api.POST("/kube/create_config_map", s.kubescontroller.CreateOrUpdateConfigMapRequest)
		api.POST("/kube/create_secret", s.kubescontroller.CreateOrUpdateSecretRequest)
		api.POST("/kube/create_namespace", s.kubescontroller.CreateNamespaceRequest)
		api.POST("/kube/create_pv", s.kubescontroller.CreatePersistentVolumeRequest)
		api.POST("/kube/create_pvc", s.kubescontroller.CreatePersistentVolumeClaimRequest)
		api.POST("/kube/create_nodeport", s.kubescontroller.CreateNodePortRequest)
		api.POST("/helm_create", s.kubescontroller.HCreateReleaseRequest)
		api.POST("helm_create_repo", s.kubescontroller.HCreateRepoRequest)
		api.POST("/kube/create_role", s.kubescontroller.CreateRoleRequest)
		api.POST("/kube/role_bind", s.kubescontroller.CreateRoleBindingRequest)
		api.POST("/kube/create_account", s.kubescontroller.CreateServiceAccountRequest)
		api.DELETE("/workspace/:id", s.TestController.DeleteWorkspace)
		api.DELETE("workspace/:id/:node_id", s.TestController.DeleteNode)
		api.DELETE("/kube_delete/:namespace/:pod_name", s.kubescontroller.DeletePodRequest)
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
	kubescontroller kubescontrollers.KubeController,
) TestRoutes {
	return TestRoutes{
		handler:         handler,
		logger:          logger,
		TestController:  TestController,
		kubescontroller: kubescontroller,
	}
}
