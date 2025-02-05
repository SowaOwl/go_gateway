package gateway

import (
	"bytes"
	"errors"
	"fmt"
	"gateway/app/gateway/DTOs"
	"gateway/app/http"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

type Service interface {
	Get(dto DTOs.GetDTO) (DTOs.ResponseDTO, error)
	Post(dto DTOs.PostDTO) (DTOs.ResponseDTO, error)
	WithBody(dto DTOs.WithBodyDTO) (DTOs.ResponseDTO, error)
}

type ServiceImpl struct {
	httpClient http.HTTPRepository
}

func NewService(repository http.HTTPRepository) Service {
	return &ServiceImpl{
		httpClient: repository,
	}
}

func (s *ServiceImpl) Get(dto DTOs.GetDTO) (DTOs.ResponseDTO, error) {
	serviceUrl, err := getServiceUrl(dto.Service)
	if err != nil {
		return DTOs.ResponseDTO{}, err
	}

	requestUrl := fmt.Sprintf("%s/api%s", serviceUrl, dto.Route)
	if dto.Params != "" {
		requestUrl = fmt.Sprintf("%s?%s", requestUrl, dto.Params)
	}

	response, body, err := s.httpClient.Get(requestUrl, dto.Bearer)
	if err != nil {
		return DTOs.ResponseDTO{}, err
	}

	responseDTO := DTOs.ResponseDTO{
		Body:        body,
		ContentType: response.Header.Get("Content-Type"),
		Status:      response.StatusCode,
	}

	return responseDTO, err
}

func (s *ServiceImpl) Post(dto DTOs.PostDTO) (DTOs.ResponseDTO, error) {
	serviceUrl, err := getServiceUrl(dto.Service)
	if err != nil {
		return DTOs.ResponseDTO{}, err
	}

	requestUrl := fmt.Sprintf("%s/api%s", serviceUrl, dto.Route)
	if dto.UrlParams != "" {
		requestUrl = fmt.Sprintf("%s?%s", requestUrl, dto.UrlParams)
	}

	reqBody := &bytes.Buffer{}

	switch dto.ContentType {
	case "multipart/form-data":
		contentType := ""
		reqBody, contentType, err = processFormData(dto.Context)
		if err != nil {
			return DTOs.ResponseDTO{}, err
		}
		dto.ContentType = contentType
	default:
		reqBody, err = processJson(dto.Context)
		if err != nil {
			return DTOs.ResponseDTO{}, err
		}
	}

	response, body, err := s.httpClient.Post(requestUrl, dto.Bearer, dto.ContentType, reqBody)
	if err != nil {
		return DTOs.ResponseDTO{}, err
	}

	responseDTO := DTOs.ResponseDTO{
		Body:        body,
		ContentType: response.Header.Get("Content-Type"),
		Status:      response.StatusCode,
	}

	return responseDTO, err
}

func (s *ServiceImpl) WithBody(dto DTOs.WithBodyDTO) (DTOs.ResponseDTO, error) {
	serviceUrl, err := getServiceUrl(dto.Service)
	if err != nil {
		return DTOs.ResponseDTO{}, err
	}

	requestUrl := fmt.Sprintf("%s/api%s", serviceUrl, dto.Route)
	if dto.UrlParams != "" {
		requestUrl = fmt.Sprintf("%s?%s", requestUrl, dto.UrlParams)
	}

	reqBody, err := processJson(dto.Context)
	if err != nil {
		return DTOs.ResponseDTO{}, err
	}

	response, body, err := s.httpClient.WithBody(requestUrl, dto.Bearer, dto.Type, reqBody)
	if err != nil {
		return DTOs.ResponseDTO{}, err
	}

	responseDTO := DTOs.ResponseDTO{
		Body:        body,
		ContentType: response.Header.Get("Content-Type"),
		Status:      response.StatusCode,
	}

	return responseDTO, err
}

func getServiceUrl(serviceName string) (string, error) {
	envRoute := fmt.Sprintf("SERVICE_%s_URL", strings.ToUpper(serviceName))

	serviceUrl := os.Getenv(envRoute)
	if serviceUrl == "" {
		return "", errors.New(fmt.Sprintf("service %s not exists", serviceName))
	}

	return serviceUrl, nil
}

func processJson(context *gin.Context) (*bytes.Buffer, error) {
	rawData, err := io.ReadAll(context.Request.Body)
	if err != nil {
		return &bytes.Buffer{}, err
	}

	return bytes.NewBuffer(rawData), nil
}

func processFormData(context *gin.Context) (*bytes.Buffer, string, error) {
	form, err := context.MultipartForm()
	if err != nil {
		return &bytes.Buffer{}, "", err
	}

	files := form.File
	params := context.Request.PostForm

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, fileHeaders := range files {
		for _, fileHeader := range fileHeaders {
			part, err := writer.CreateFormFile(key, fileHeader.Filename)
			if err != nil {
				return &bytes.Buffer{}, "", err
			}

			if err = func() error {
				file, err := fileHeader.Open()
				if err != nil {
					return err
				}

				var returnErr error
				defer func(file multipart.File) {
					err := file.Close()
					if err != nil {
						returnErr = err
					}
				}(file)

				_, err = io.Copy(part, file)
				if err != nil {
					return err
				}

				return returnErr
			}(); err != nil {
				return &bytes.Buffer{}, "", err
			}
		}
	}

	for key, values := range params {
		for _, value := range values {
			if err := writer.WriteField(key, value); err != nil {
				return &bytes.Buffer{}, "", err
			}
		}
	}

	err = writer.Close()
	if err != nil {
		return &bytes.Buffer{}, "", err
	}

	return body, writer.FormDataContentType(), nil
}
