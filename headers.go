package http_proxy

func (requestIntent *ProxiedRequestImpl) AddHeader(key string, value string) ProxiedRequest {
	requestIntent.verifyUnderlyingRequestNotGenerated()
	if _, isFound := requestIntent.headers[key]; !isFound {
		requestIntent.headers[key] = []string{}
	}
	requestIntent.headers[key] = append(requestIntent.headers[key], value)
	return requestIntent
}

func (requestIntent *ProxiedRequestImpl) SetHeader(key string, value string) ProxiedRequest {
	requestIntent.verifyUnderlyingRequestNotGenerated()
	requestIntent.headers[key] = []string{value}
	return requestIntent
}

func (requestIntent *ProxiedRequestImpl) SetHeaders(headers map[string]string) ProxiedRequest {
	requestIntent.verifyUnderlyingRequestNotGenerated()
	for key, value := range headers {
		requestIntent.headers[key] = []string{value}
	}
	return requestIntent
}

func (requestIntent *ProxiedRequestImpl) SetMultiValueHeaders(headers map[string][]string) ProxiedRequest {
	requestIntent.verifyUnderlyingRequestNotGenerated()
	for key, values := range headers {
		requestIntent.headers[key] = values
	}
	return requestIntent
}
