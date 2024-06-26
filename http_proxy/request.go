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
	// Adds an interceptor that is executed over the response
	WithGenericInterceptor(handlers ...errorHandler) ProxiedRequest
	// Adds an interceptor that is executed when the response status code
	// matches the provided value
	WithStatusCodeInterceptor(statusCode int, handlers ...errorHandler) ProxiedRequest
	// Generates the underlying request without sending it. After this the request
	// can't be modified or it will return an error
	UnderlyingRequest() (*http.Request, error)
	// Set the context of the request
	WithContext(ctx context.Context) ProxiedRequest
	// Generates the underlying request if not already generated and sends it
	Send() (*http.Response, error)
}

type proxiedRequestImpl struct {
	method                 string
	url                    string
	body                   io.Reader
	context                context.Context
	headers                map[string][]string
	requestError           error
	underlyingRequest      *http.Request
	statusCodeInterceptors map[int][]errorHandler
	genericInterceptors    []errorHandler
}

func NewRequest(method string, url string) *proxiedRequestImpl {
	return &proxiedRequestImpl{
		method:                 method,
		headers:                map[string][]string{},
		url:                    url,
		body:                   http.NoBody,
		genericInterceptors:    []errorHandler{},
		statusCodeInterceptors: map[int][]errorHandler{},
	}
}

func (requestIntent *proxiedRequestImpl) UnderlyingRequest() (*http.Request, error) {
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

func (requestIntent *proxiedRequestImpl) Send() (*http.Response, error) {
	if requestIntent.underlyingRequest == nil {
		requestIntent.UnderlyingRequest()
	}

	if requestIntent.requestError != nil {
		return nil, requestIntent.requestError
	}

	var response *http.Response
	var err error
	if response, err = http.DefaultClient.Do(requestIntent.underlyingRequest); err == nil {
		response, err = requestIntent.validateResponse(response)
	}
	return response, err
}

func (requestIntent *proxiedRequestImpl) verifyUnderlyingRequestNotGenerated() {
	if requestIntent.underlyingRequest != nil {
		requestIntent.requestError = fmt.Errorf("tried to modify request proxy after generating the underlying request")
	}
}
