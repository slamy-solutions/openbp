package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Parameter struct {
	name  string
	value interface{}
}

type Header struct {
	name  string
	value string
}

func formatURL(path string, params []*Parameter) string {
	baseURL := "http://example.com"
	u, _ := url.ParseRequestURI(baseURL)
	u.Path = path

	requestParams := url.Values{}
	for _, prm := range params {
		requestParams.Add(prm.name, fmt.Sprintf("%v", prm.value))
	}
	u.RawQuery = requestParams.Encode()

	return fmt.Sprintf("%v", u)
}

func baseRequest(path string, params []*Parameter, headers []*Header, method string, boby io.Reader) ([]byte, int, error) {
	fullURL := formatURL(path, params)

	req, _ := http.NewRequest(method, fmt.Sprintf("%v", fullURL), boby)
	req.Header = http.Header{}
	for _, header := range headers {
		req.Header.Add(header.name, header.value)
	}
	client := http.Client{}
	defer client.CloseIdleConnections()

	response, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, 0, err
	}
	response.Body.Close()

	return body, response.StatusCode, nil
}

func PostJSON(path string, params []*Parameter, headers []*Header, in interface{}, out interface{}) (int, error) {
	requestData, err := json.Marshal(&in)
	if err != nil {
		return 0, err
	}
	reader := bytes.NewReader(requestData)

	responseData, code, err := baseRequest(path, params, headers, "POST", reader)
	if err != nil {
		return 0, err
	}

	return code, json.Unmarshal(responseData, &out)
}

func PatchJSON(path string, params []*Parameter, headers []*Header, in interface{}, out interface{}) (int, error) {
	requestData, err := json.Marshal(&in)
	if err != nil {
		return 0, err
	}
	reader := bytes.NewReader(requestData)

	responseData, code, err := baseRequest(path, params, headers, "PATCH", reader)
	if err != nil {
		return 0, err
	}

	return code, json.Unmarshal(responseData, &out)
}

func GetJSON(path string, params []*Parameter, headers []*Header, out interface{}) (int, error) {
	data, code, err := Get(path, params, headers)
	if err != nil {
		return 0, err
	}

	return code, json.Unmarshal(data, &out)
}

func Get(path string, params []*Parameter, headers []*Header) ([]byte, int, error) {
	return baseRequest(path, params, headers, "GET", nil)
}

func DeleteJSON(path string, params []*Parameter, headers []*Header, out interface{}) (int, error) {
	data, code, err := Delete(path, params, headers)
	if err != nil {
		return 0, err
	}

	return code, json.Unmarshal(data, &out)
}

func Delete(path string, params []*Parameter, headers []*Header) ([]byte, int, error) {
	return baseRequest(path, params, headers, "DELETE", nil)
}
