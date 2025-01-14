package DTOs

import "github.com/gin-gonic/gin"

type WithBodyDTO struct {
	Service   string
	Route     string
	UrlParams string
	Bearer    string
	Context   *gin.Context
	Type      string
}
