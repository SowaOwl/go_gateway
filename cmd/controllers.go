package cmd

import "gateway/app/gateway"

func InitControllers() gateway.Controller {
	// Repositories
	gatewayRepository := gateway.NewHTTPRepository()

	// Services
	gatewayService := gateway.NewService(gatewayRepository)

	// Controllers
	gatewayController := gateway.NewController(gatewayService)

	return gatewayController
}
