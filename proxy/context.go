package http_proxy

import (
	"context"
)

func (requestIntent *proxiedRequestImpl) WithContext(ctx context.Context) ProxiedRequest {
	requestIntent.context = ctx
	return requestIntent
}
