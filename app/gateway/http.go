package gateway

import (
	"io"
	"net/http"
)

type HTTPRepository interface {
	Get(url string, bearerToken string) (*http.Response, []byte, error)
	Post(url string, bearerToken string, contentType string, body io.Reader) (*http.Response, []byte, error)
	Delete(url string, bearerToken string, body io.Reader) (*http.Response, []byte, error)
}

type HTTPRepositoryImpl struct {
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

func NewHTTPRepository() HTTPRepository {
	return &HTTPRepositoryImpl{
		client: http.Client{},
	}
}

func (h *HTTPRepositoryImpl) Get(url string, bearerToken string) (*http.Response, []byte, error) {
	h.client.Transport = &TransportWithToken{
		Token: bearerToken,
	}

	resp, err := h.client.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	return resp, body, nil
}

func (h *HTTPRepositoryImpl) Post(url string, bearerToken string, contentType string, reqData io.Reader) (*http.Response, []byte, error) {
	h.client.Transport = &TransportWithToken{
		Token: bearerToken,
	}

	resp, err := h.client.Post(url, contentType, reqData)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	return resp, body, nil
}

func (h *HTTPRepositoryImpl) Delete(url string, bearerToken string, reqData io.Reader) (*http.Response, []byte, error) {
	h.client.Transport = &TransportWithToken{
		Token: bearerToken,
	}

	req, err := http.NewRequest(http.MethodDelete, url, reqData)
	if err != nil {
		return nil, nil, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	return resp, body, nil
}
