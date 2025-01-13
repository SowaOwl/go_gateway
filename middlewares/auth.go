package middlewares

import (
	"gateway/app/jwt"
	"gateway/app/redis"
	"gateway/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func BearerTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			util.SendError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.SendError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		tokenString := parts[1]

		jwtService, err := jwt.NewJwt()
		if err != nil {
			util.SendError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		id, err := jwtService.GetIdFromToken(tokenString)
		if err != nil {
			util.SendError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		redisService, err := redis.NewRedisService()
		if err != nil {
			util.SendError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		redisData, err := redisService.GetClient(id)
		if err != nil {
			util.SendError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if redisData == "" {
			util.SendError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		c.Next()
	}
}
