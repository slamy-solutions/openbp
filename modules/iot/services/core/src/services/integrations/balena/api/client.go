package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	url "net/url"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var ErrServerURLInvalid = errors.New("server URL has invalid format")
var ErrFailedToCreateRequest = errors.New("failed to create request object")
var ErrServerConnectionError = errors.New("failed to connect to the balena server")
var ErrUnauthenticated = errors.New("balena responded with 401 status code (unauthentiocated). Probably auth token is invalid")
var ErrUnauthorized = errors.New("balena responded with 403 status code (unauthorized). Probably you dont have permissions to perform this action")
var ErrInternalBalenaError = errors.New("balena responded with >=500 status code (server error)")
var ErrInvalidStatusCode = errors.New("balena responded with invalid status code")
var ErrFailedToParseResponse = errors.New("failed to parse response body from the balena server")

type client struct {
	httpClient *http.Client
}

type Client interface {
	GetAllDevices(ctx context.Context, server BalenaServerInfo) ([]Device, error)
	Ping(ctx context.Context, server BalenaServerInfo) error

	// Close client and free resources
	Close()
}

func NewClient() Client {
	tr := otelhttp.NewTransport(&http.Transport{
		MaxIdleConns:          10,
		IdleConnTimeout:       15 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		DisableKeepAlives:     false,
	})

	return &client{
		httpClient: &http.Client{
			Transport: tr,
		},
	}
}

func (c *client) Close() {
	c.httpClient.CloseIdleConnections()
}

func (c *client) getErrorByResponse(response *http.Response) error {
	if response.StatusCode == 401 {
		return errors.Join(ErrUnauthenticated, ErrInvalidStatusCode)
	}
	if response.StatusCode == 403 {
		return errors.Join(ErrUnauthorized, ErrInvalidStatusCode)
	}
	if response.StatusCode >= 500 {
		return errors.Join(ErrInternalBalenaError, ErrInvalidStatusCode)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return ErrInvalidStatusCode
	}

	return nil
}

func (c *client) GetAllDevices(ctx context.Context, server BalenaServerInfo) ([]Device, error) {
	requestURL := fmt.Sprintf("%s/v6/device", server.BaseURL)
	_, err := url.ParseRequestURI(requestURL)
	if err != nil {
		return nil, errors.Join(ErrServerURLInvalid, err)
	}

	r, err := http.NewRequestWithContext(ctx, "GET", requestURL, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, errors.Join(ErrFailedToCreateRequest, err)
	}
	r.Header["Content-Type"] = []string{"application/json"}
	r.Header["Authorization"] = []string{"Bearer " + server.APIToken}

	response, err := c.httpClient.Do(r)
	if err != nil {
		return nil, errors.Join(ErrServerConnectionError, err)
	}
	defer response.Body.Close()

	if err = c.getErrorByResponse(response); err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Join(ErrServerConnectionError, err)
	}

	var devicesResponse struct {
		Devices []Device `json:"d"`
	}
	err = json.Unmarshal(body, &devicesResponse)
	if err != nil {
		return nil, errors.Join(ErrFailedToParseResponse, err)
	}

	return devicesResponse.Devices, nil
}

func (c *client) Ping(ctx context.Context, server BalenaServerInfo) error {
	requestURL := fmt.Sprintf("%s/user/v1/whoami", server.BaseURL)
	_, err := url.ParseRequestURI(requestURL)
	if err != nil {
		return errors.Join(ErrServerURLInvalid, err)
	}

	r, err := http.NewRequestWithContext(ctx, "GET", requestURL, bytes.NewBuffer([]byte{}))
	if err != nil {
		return errors.Join(ErrFailedToCreateRequest, err)
	}
	r.Header["Content-Type"] = []string{"application/json"}
	r.Header["Authorization"] = []string{"Bearer " + server.APIToken}

	response, err := c.httpClient.Do(r)
	if err != nil {
		return errors.Join(ErrServerConnectionError, err)
	}
	defer response.Body.Close()

	if err = c.getErrorByResponse(response); err != nil {
		return err
	}

	return nil
}
