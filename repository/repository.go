package repository

import "go.uber.org/fx"

// Module exports dependency
var Module = fx.Options(
	fx.Provide(NewTestRepository),
	fx.Provide(NewKubernetesRepository),
)
