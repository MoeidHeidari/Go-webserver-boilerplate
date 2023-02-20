package bootstrap

import (
	"main/api/controllers"
	"main/api/kubescontrollers"
	"main/api/middlewares"
	"main/api/routes"
	"main/api/ws"
	"main/lib"
	"main/repository"
	"main/services"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	controllers.Module,
	routes.Module,
	lib.Module,
	services.Module,
	middlewares.Module,
	repository.Module,
	kubescontrollers.Module,
	ws.Module,
)
