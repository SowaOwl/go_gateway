package gateway

import (
	"errors"
	"fmt"
	"gateway/app/gateway/DTOs"
	"os"
	"strings"
)

type Service interface {
	Get(dto DTOs.GetDTO) ([]byte, string, error)
}

type ServiceImpl struct {
	httpClient HTTPRepository
}

func NewService(repository HTTPRepository) *ServiceImpl {
	return &ServiceImpl{
		httpClient: repository,
	}
}

func (s *ServiceImpl) Get(dto DTOs.GetDTO) ([]byte, string, error) {
	serviceUrl, err := getServiceUrl(dto.Service)
	if err != nil {
		return nil, "", err
	}

	requestUrl := fmt.Sprintf("%s/api%s", serviceUrl, dto.Route)
	if dto.Params != "" {
		requestUrl = fmt.Sprintf("%s?%s", requestUrl, dto.Params)
	}

	response, body, err := s.httpClient.Get(requestUrl, dto.Bearer)

	return body, response.Header.Get("Content-Type"), err
}

func getServiceUrl(serviceName string) (string, error) {
	envRoute := fmt.Sprintf("SERVICE_%s_URL", strings.ToUpper(serviceName))

	serviceUrl := os.Getenv(envRoute)
	if serviceUrl == "" {
		return "", errors.New(fmt.Sprintf("service %s not exists", serviceName))
	}

	return serviceUrl, nil
}
