package gateway

import (
	"gateway/app/gateway/DTOs"
	"gateway/util"
	"github.com/gin-gonic/gin"
	"net/http"
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
		util.SendError(context, http.StatusInternalServerError, err.Error(), "")
		return
	}

	dto := DTOs.GetDTO{
		Service: urlParams.Service,
		Route:   urlParams.Route,
		Params:  context.Request.URL.RawQuery,
		Bearer:  context.GetHeader("Authorization"),
	}

	response, err := c.service.Get(dto)
	if err != nil {
		util.SendError(context, http.StatusInternalServerError, err.Error(), "")
		return
	}

	context.Data(response.Status, response.ContentType, response.Body)
	return
}

func (c *ControllerImpl) Post(context *gin.Context) {
	var urlParams UrlParam
	if err := context.ShouldBindUri(&urlParams); err != nil {
		util.SendError(context, http.StatusInternalServerError, err.Error(), "")
		return
	}

	dto := DTOs.PostDTO{
		Service:     urlParams.Service,
		Route:       urlParams.Route,
		UrlParams:   context.Request.URL.RawQuery,
		Bearer:      context.GetHeader("Authorization"),
		ContentType: context.ContentType(),
		Context:     context,
	}

	response, err := c.service.Post(dto)
	if err != nil {
		util.SendError(context, http.StatusInternalServerError, err.Error(), "")
		return
	}

	context.Data(response.Status, response.ContentType, response.Body)
	return
}

func (c *ControllerImpl) WithBody(context *gin.Context) {
	var urlParams UrlParam
	if err := context.ShouldBindUri(&urlParams); err != nil {
		util.SendError(context, http.StatusInternalServerError, err.Error(), "")
		return
	}

	dto := DTOs.WithBodyDTO{
		Service:   urlParams.Service,
		Route:     urlParams.Route,
		UrlParams: context.Request.URL.RawQuery,
		Bearer:    context.GetHeader("Authorization"),
		Context:   context,
		Type:      context.Request.Method,
	}

	response, err := c.service.WithBody(dto)
	if err != nil {
		util.SendError(context, http.StatusInternalServerError, err.Error(), "")
		return
	}

	context.Data(response.Status, response.ContentType, response.Body)
	return
}
