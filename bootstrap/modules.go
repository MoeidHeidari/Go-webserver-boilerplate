package bootstrap

import (
	"go.uber.org/fx"
	"main/api/controllers"
	"main/api/middlewares"
	"main/api/routes"
	"main/lib"
	"main/repository"
	"main/services"
)

var CommonModules = fx.Options(
	controllers.Module,
	routes.Module,
	lib.Module,
	services.Module,
	middlewares.Module,
	repository.Module,
)
