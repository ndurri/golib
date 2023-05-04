package service

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ServiceConfig struct {
	Headers    map[string]string
	Operations map[string]string
}

type NVP map[string]string

func (params NVP) ToURLValues() url.Values {
	vals := url.Values{}
	for key, value := range params {
		vals.Add(key, value)
	}
	return vals
}

func BodyParams(params NVP) *bytes.Buffer {
	uv := params.ToURLValues()
	return bytes.NewBufferString(uv.Encode())
}

func parsePathParams(endpoint string, params NVP) string {
	parsed := endpoint
	for name, value := range params {
		parsed = strings.ReplaceAll(parsed, ":"+name, value)
	}
	return parsed
}

func parseURL(endpoint string, pathParams NVP, urlParams NVP) (*url.URL, error) {
	var parsed string
	if pathParams == nil {
		parsed = endpoint
	} else {
		parsed = parsePathParams(endpoint, pathParams)
	}
	serviceURL, err := url.Parse(parsed)
	if err != nil {
		return nil, err
	}
	if urlParams != nil {
		serviceURL.RawQuery = urlParams.ToURLValues().Encode()
	}
	return serviceURL, nil
}

var (
	AuthError       = errors.New("401 Not Authorized")
	NotFoundError   = errors.New("404 Not Found")
	PayloadError    = errors.New("4xx Bad Request")
	ServerError     = errors.New("5xx Server Error")
	UnexpectedError = errors.New("Unexpected Response")
)

func Get(endpoint string, headers NVP, pathParams NVP, urlParams NVP) ([]byte, error) {
	return call(http.MethodGet, endpoint, headers, pathParams, urlParams, nil)
}

func Patch(endpoint string, headers NVP, pathParams NVP, urlParams NVP, body io.Reader) ([]byte, error) {
	return call(http.MethodPatch, endpoint, headers, pathParams, urlParams, body)
}

func Post(endpoint string, headers NVP, pathParams NVP, urlParams NVP, body io.Reader) ([]byte, error) {
	return call(http.MethodPost, endpoint, headers, pathParams, urlParams, body)
}

func call(method string, endpoint string, headers NVP, pathParams NVP, urlParams NVP, body io.Reader) ([]byte, error) {
	serviceURL, err := parseURL(endpoint, pathParams, urlParams)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, serviceURL.String(), body)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	resbody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return resbody, nil
	} else {
		if res.StatusCode == 401 {
			return nil, AuthError
		} else if res.StatusCode == 404 {
			return nil, NotFoundError
		} else if res.StatusCode >= 400 && res.StatusCode <= 499 {
			return nil, PayloadError
		} else if res.StatusCode >= 500 && res.StatusCode <= 599 {
			return nil, ServerError
		} else {
			return nil, UnexpectedError
		}
	}
}
