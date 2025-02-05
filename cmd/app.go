package cmd

import (
	"gateway/app/gateway"
	"gateway/app/http"
	"gateway/middlewares"
)

type App struct {
	Dependencies *Dependencies
	Controllers  Controllers
	Middlewares  Middlewares
}

type Controllers struct {
	Gateway gateway.Controller
}

type Middlewares struct {
	Auth *middlewares.AuthMiddleware
	Log  *middlewares.LogMiddleware
}

func InitApp() *App {
	deps := InitDependencies()

	//Repositories
	httpRepository := http.NewHTTPRepository()

	// Services
	gatewayService := gateway.NewService(httpRepository)

	// Controllers
	gatewayController := gateway.NewController(gatewayService)

	controllers := Controllers{
		Gateway: gatewayController,
	}

	middlewaresObj := Middlewares{
		Auth: middlewares.NewAuthMiddleware(deps.DB, deps.Jwt, deps.Redis),
		Log:  middlewares.NewLogMiddleware(deps.DB),
	}

	return &App{
		Dependencies: deps,
		Controllers:  controllers,
		Middlewares:  middlewaresObj,
	}
}
