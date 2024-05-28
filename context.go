package http_proxy

import (
	"context"
)

func (requestIntent *ProxiedRequestImpl) WithContext(ctx context.Context) ProxiedRequest {
	requestIntent.context = ctx
	return requestIntent
}
