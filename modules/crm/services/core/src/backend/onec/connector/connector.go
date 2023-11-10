package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
)

type OneCConnector struct {
	Client *http.Client

	ServerURL   string
	ServerToken string
}

func NewOneCConnector(serverURL string, serverToken string) *OneCConnector {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.IdleConnTimeout = 90 * time.Second
	transport.MaxIdleConnsPerHost = 2
	transport.MaxConnsPerHost = 10

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	return &OneCConnector{
		Client:      httpClient,
		ServerURL:   serverURL,
		ServerToken: serverToken,
	}
}

var ErrRequestFailed = errors.New("request failed")

func MakeRequest[R any](ctx context.Context, connector *OneCConnector, method string, targetURL string, reqeustData interface{}) (*R, int, error) {
	var requestDataBytes []byte
	var err error

	if reqeustData != nil {
		requestDataBytes, err = json.Marshal(reqeustData)
		if err != nil {
			err := errors.Join(errors.New("failed to marshal request data"), err)
			return nil, 0, err
		}
	}

	request, err := http.NewRequestWithContext(ctx, method, targetURL, bytes.NewBuffer(requestDataBytes))
	if err != nil {
		err := errors.Join(errors.New("failed to create request"), err)
		return nil, 0, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := connector.Client.Do(request)
	if err != nil {
		err := errors.Join(models.ErrBackendUnavailable, ErrRequestFailed, err)
		return nil, 0, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode, nil
	}

	var responseData R
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New("failed to unmarshal response"), err)
		return nil, 0, err
	}

	return &responseData, response.StatusCode, nil
}
