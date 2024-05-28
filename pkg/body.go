package http_proxy

import (
	"bytes"
	"encoding/json"
	"io"
)

func (requestIntent *ProxiedRequestImpl) SetBody(body io.Reader) ProxiedRequest {
	requestIntent.verifyUnderlyingRequestNotGenerated()
	requestIntent.body = body
	return requestIntent
}

func (requestIntent *ProxiedRequestImpl) SetJSONBody(body any) ProxiedRequest {
	requestIntent.verifyUnderlyingRequestNotGenerated()
	payload, marshalErr := json.Marshal(body)
	requestIntent.requestError = marshalErr
	if marshalErr != nil {
		return requestIntent
	}
	return requestIntent.SetBody(bytes.NewBuffer(payload)).SetHeader("Content-Type", "application/json")
}
