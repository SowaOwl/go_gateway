package gateway

import (
	"gateway/app/gateway/DTOs"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Get(context *gin.Context)
	Post(context *gin.Context)
	WithBody(context *gin.Context)
}

type ControllerImpl struct {
	service Service
}

func NewController(service Service) Controller {
	return &ControllerImpl{
		service: service,
	}
}

type UrlParam struct {
	Service string `uri:"service" binding:"required"`
	Route   string `uri:"route"`
}

func (c *ControllerImpl) Get(context *gin.Context) {
	var urlParams UrlParam
	if err := context.ShouldBindUri(&urlParams); err != nil {
		context.JSON(400, gin.H{"msg": err.Error()})
	}

	token := context.GetHeader("Authorization")

	dto := DTOs.GetDTO{
		Service: urlParams.Service,
		Route:   urlParams.Route,
		Params:  context.Request.URL.RawQuery,
		Bearer:  token,
	}

	response, err := c.service.Get(dto)
	if err != nil {
		context.JSON(400, gin.H{"msg": err.Error()})
	}

	context.Data(response.Status, response.ContentType, response.Body)
	return
}

func (c *ControllerImpl) Post(context *gin.Context) {
	var urlParams UrlParam
	if err := context.ShouldBindUri(&urlParams); err != nil {
		context.JSON(400, gin.H{"msg": err.Error()})
	}

	token := context.GetHeader("Authorization")

	dto := DTOs.PostDTO{
		Service:     urlParams.Service,
		Route:       urlParams.Route,
		UrlParams:   context.Request.URL.RawQuery,
		Bearer:      token,
		ContentType: context.ContentType(),
		Context:     context,
	}

	response, err := c.service.Post(dto)
	if err != nil {
		context.JSON(400, gin.H{"msg": err.Error()})
	}

	context.Data(response.Status, response.ContentType, response.Body)
	return
}

func (c *ControllerImpl) WithBody(context *gin.Context) {
	var urlParams UrlParam
	if err := context.ShouldBindUri(&urlParams); err != nil {
		context.JSON(400, gin.H{"msg": err.Error()})
	}

	token := context.GetHeader("Authorization")

	dto := DTOs.WithBodyDTO{
		Service:   urlParams.Service,
		Route:     urlParams.Route,
		UrlParams: context.Request.URL.RawQuery,
		Bearer:    token,
		Context:   context,
		Type:      context.Request.Method,
	}

	response, err := c.service.WithBody(dto)
	if err != nil {
		context.JSON(400, gin.H{"msg": err.Error()})
	}

	context.Data(response.Status, response.ContentType, response.Body)
	return
}
