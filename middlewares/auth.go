package middlewares

import (
	"errors"
	"fmt"
	"gateway/app/jwt"
	"gateway/app/redis"
	"gateway/database/model"
	"gateway/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	db           *gorm.DB
	jwtService   jwt.Service
	redisService redis.Service
}

func NewAuthMiddleware(db *gorm.DB, jwtService jwt.Service, redisService redis.Service) *AuthMiddleware {
	return &AuthMiddleware{
		db:           db,
		jwtService:   jwtService,
		redisService: redisService,
	}
}

func (a *AuthMiddleware) BearerTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		skipAuth, err := a.ifRouteAcceptToSkip(c.Request.URL.Path)
		if err != nil {
			a.processError(err, c)
			return
		}

		if skipAuth {
			c.Next()
			return
		}

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			a.processError(fmt.Errorf("authorization field is empty"), c)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			a.processError(fmt.Errorf("no token found in request"), c)
			return
		}

		tokenString := parts[1]

		id, err := a.jwtService.GetIdFromToken(tokenString)
		if err != nil {
			a.processError(err, c)
			return
		}

		redisData, err := a.redisService.GetClient(id)
		if err != nil {
			a.processError(err, c)
			return
		}

		if redisData == "" {
			a.processError(fmt.Errorf("user data in redis is empty"), c)
			return
		}

		c.Next()
	}
}

func (a *AuthMiddleware) processError(err error, c *gin.Context) {
	util.SaveErrToDB(err, a.db)
	util.SendError(c, http.StatusUnauthorized, "Unauthorized", "")
	return
}

func (a *AuthMiddleware) ifRouteAcceptToSkip(route string) (bool, error) {
	var endpoint model.WithoutAuthEndpoint

	if err := a.db.Where("value = ?", route).First(&endpoint).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}
