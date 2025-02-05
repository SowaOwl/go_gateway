package http

import (
	"io"
	"net/http"
	"time"
)

const defaultTimeout = 300 * time.Second

type Repository interface {
	Get(url string, bearerToken string) (*http.Response, []byte, error)
	Post(url string, bearerToken string, contentType string, reqData io.Reader) (*http.Response, []byte, error)
	WithBody(url string, bearerToken string, reqType string, reqData io.Reader) (*http.Response, []byte, error)
}

type RepositoryImpl struct {
	client http.Client
}

type TransportWithToken struct {
	Token     string
	Transport http.RoundTripper
}

func (t *TransportWithToken) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.Token)
	if t.Transport != nil {
		return t.Transport.RoundTrip(req)
	}

	return http.DefaultTransport.RoundTrip(req)
}

func NewHTTPRepository() Repository {
	return &RepositoryImpl{
		client: http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func (h *RepositoryImpl) Get(url string, bearerToken string) (*http.Response, []byte, error) {
	h.client.Transport = &TransportWithToken{
		Token: bearerToken,
	}

	resp, err := h.client.Get(url)
	if err != nil {
		return nil, nil, err
	}

	var returnErr error
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			returnErr = err
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)

	return resp, body, returnErr
}

func (h *RepositoryImpl) Post(url string, bearerToken string, contentType string, reqData io.Reader) (*http.Response, []byte, error) {
	h.client.Transport = &TransportWithToken{
		Token: bearerToken,
	}

	resp, err := h.client.Post(url, contentType, reqData)
	if err != nil {
		return nil, nil, err
	}

	var returnErr error
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			returnErr = err
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)

	return resp, body, returnErr
}

func (h *RepositoryImpl) WithBody(url string, bearerToken string, reqType string, reqData io.Reader) (*http.Response, []byte, error) {
	h.client.Transport = &TransportWithToken{
		Token: bearerToken,
	}

	req, err := http.NewRequest(reqType, url, reqData)
	if err != nil {
		return nil, nil, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	var returnErr error
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			returnErr = err
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)

	return resp, body, returnErr
}
