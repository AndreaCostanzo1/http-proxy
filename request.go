package http_proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type ProxiedRequest interface {
	// Adds the key-value pair to the header.
	// It appends to any existing values associated with key
	AddHeader(key string, value string) ProxiedRequest
	// Applies the body to request replacing the older one if present
	SetBody(body io.Reader) ProxiedRequest
	// Set the key-value pair to the header.
	// It replaces any existing values associated with key
	SetHeader(key string, value string) ProxiedRequest
	// Set the key-value pairs to the header.
	// It replaces any existing values associated with the keys
	SetHeaders(headers map[string]string) ProxiedRequest
	// Transform the passed object into an io.Reader and applies it as body
	// of the request, replacing the older one if present
	SetJSONBody(body any) ProxiedRequest
	// It allows to set comma separated values for the provided keys
	// It replaces any existing values associated with the keys
	SetMultiValueHeaders(headers map[string][]string) ProxiedRequest
	// Generates the underlying request without sending it. After this the request
	// can't be modified or it will return an error
	UnderlyingRequest() (*http.Request, error)
	// Generates the underlying request if not already generated and sends it
	Send() (*http.Response, error)
}

type ProxiedRequestImpl struct {
	method            string
	url               string
	body              io.Reader
	context           context.Context
	headers           map[string][]string
	requestError      error
	underlyingRequest *http.Request
}

func NewRequest(method string, url string) *ProxiedRequestImpl {
	return &ProxiedRequestImpl{
		method:  method,
		headers: map[string][]string{},
		url:     url,
		body:    http.NoBody,
	}
}

func (requestIntent *ProxiedRequestImpl) UnderlyingRequest() (*http.Request, error) {
	if requestIntent.requestError != nil {
		return nil, requestIntent.requestError
	}
	if requestIntent.underlyingRequest != nil {
		return requestIntent.underlyingRequest, nil
	}
	newRequest, createRequestErr := http.NewRequest(requestIntent.method, requestIntent.url, requestIntent.body)
	requestIntent.underlyingRequest = newRequest
	requestIntent.requestError = createRequestErr
	if createRequestErr == nil {
		for headerKey, headerValues := range requestIntent.headers {
			for _, value := range headerValues {
				requestIntent.underlyingRequest.Header.Add(headerKey, value)
			}
		}
		if requestIntent.context != nil {
			requestIntent.underlyingRequest = requestIntent.underlyingRequest.WithContext(requestIntent.context)
		}
	}
	return newRequest, createRequestErr
}

func (requestIntent *ProxiedRequestImpl) Send() (*http.Response, error) {
	if requestIntent.underlyingRequest == nil {
		requestIntent.UnderlyingRequest()
	}

	if requestIntent.requestError != nil {
		return nil, requestIntent.requestError
	}

	return http.DefaultClient.Do(requestIntent.underlyingRequest)
}

func (requestIntent *ProxiedRequestImpl) verifyUnderlyingRequestNotGenerated() {
	if requestIntent.underlyingRequest != nil {
		requestIntent.requestError = fmt.Errorf("tried to modify request proxy after generating the underlying request")
	}
}
