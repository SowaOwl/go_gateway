package middlewares

import (
	"errors"
	"gateway/app/jwt"
	"gateway/app/redis"
	"gateway/database/model"
	"gateway/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func BearerTokenMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		skipAuth, err := ifRouteAcceptToSkip(db, c.Request.URL.Path)
		if err != nil {
			util.SendError(c, http.StatusUnauthorized, "Unauthorized", "")
			return
		}

		if skipAuth {
			c.Next()
			return
		}

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			util.SendError(c, http.StatusUnauthorized, "Unauthorized", "")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.SendError(c, http.StatusUnauthorized, "Unauthorized", "")
			return
		}

		tokenString := parts[1]

		jwtService, err := jwt.NewJwt()
		if err != nil {
			util.SendError(c, http.StatusUnauthorized, "Unauthorized", "")
			return
		}

		id, err := jwtService.GetIdFromToken(tokenString)
		if err != nil {
			util.SendError(c, http.StatusUnauthorized, "Unauthorized", "")
			return
		}

		redisService, err := redis.NewRedisService()
		if err != nil {
			util.SendError(c, http.StatusUnauthorized, "Unauthorized", "")
			return
		}

		redisData, err := redisService.GetClient(id)
		if err != nil {
			util.SendError(c, http.StatusUnauthorized, "Unauthorized", "")
			return
		}

		if redisData == "" {
			util.SendError(c, http.StatusUnauthorized, "Unauthorized", "")
			return
		}

		c.Next()
	}
}

func ifRouteAcceptToSkip(db *gorm.DB, route string) (bool, error) {
	var endpoint model.WithoutAuthEndpoint

	if err := db.Where("value = ?", route).First(&endpoint).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}
