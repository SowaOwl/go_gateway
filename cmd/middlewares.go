package cmd

import "gateway/middlewares"

func InitMiddlewares(deps *Dependencies) (*middlewares.AuthMiddleware, *middlewares.LogMiddleware) {
	authMiddleware := middlewares.NewAuthMiddleware(deps.DB, deps.Jwt, deps.Redis)
	logMiddleware := middlewares.NewLogMiddleware(deps.DB)

	return authMiddleware, logMiddleware
}
