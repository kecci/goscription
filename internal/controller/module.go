package controller

import "go.uber.org/fx"

//Module for controller module
var Module = fx.Invoke(
	InitArticleController,
	InitDomainController,
	InitUserController,
)
