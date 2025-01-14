package DTOs

import "github.com/gin-gonic/gin"

type PostDTO struct {
	Service     string
	Route       string
	UrlParams   string
	Bearer      string
	ContentType string
	Context     *gin.Context
}
