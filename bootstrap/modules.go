package bootstrap

import (
	"main/api/controllers"
	"main/api/currencies"
	"main/api/kubes"
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
	currencies.Module,
	repository.Module,
	kubes.Module,
	ws.Module,
)
