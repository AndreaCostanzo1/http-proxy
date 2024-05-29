package http_proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type format string

const FORMAT_TYPE = "__FORMAT"

const (
	FORMAT_STRING format = "String"
	FORMAT_JSON   format = "JSON"
)

type errorHandler = func(parsedBody map[string]interface{}, response *http.Response) error

func (requestIntent *proxiedRequestImpl) WithGenericInterceptor(handlers ...errorHandler) ProxiedRequest {
	requestIntent.genericInterceptors = append(requestIntent.genericInterceptors, handlers...)
	return requestIntent
}

func (requestIntent *proxiedRequestImpl) WithStatusCodeInterceptor(statusCode int, handlers ...errorHandler) ProxiedRequest {
	existingHandlers, isFound := requestIntent.statusCodeInterceptors[statusCode]
	if !isFound {
		existingHandlers = []errorHandler{}
	}
	requestIntent.statusCodeInterceptors[statusCode] = append(existingHandlers, handlers...)
	return requestIntent
}

func (requestIntent *proxiedRequestImpl) validateResponse(response *http.Response) (*http.Response, error) {
	var err error
	responseBody := extractResponseBody(response)
	if statusCodeInterceptors, isDefined := requestIntent.statusCodeInterceptors[response.StatusCode]; isDefined {
		for _, interceptor := range statusCodeInterceptors {
			err = interceptor(responseBody, response)
			if err != nil {
				return response, err
			}
		}
	}
	for _, handler := range requestIntent.genericInterceptors {
		err = handler(responseBody, response)
		if err != nil {
			return response, err
		}
	}
	return response, err
}

func extractResponseBody(response *http.Response) map[string]interface{} {
	var jsonErr error
	var responseBody map[string]interface{}
	var buf bytes.Buffer
	teeReader := io.TeeReader(response.Body, &buf)
	jsonErr = json.NewDecoder(teeReader).Decode(&responseBody)
	response.Body = io.NopCloser(&buf)
	switch {
	case jsonErr == nil:
		responseBody[FORMAT_TYPE] = FORMAT_JSON
	default:
		responseBody = map[string]interface{}{FORMAT_TYPE: FORMAT_STRING}
	}
	return responseBody
}
