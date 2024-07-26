package http_proxy_test

import (
	"testing"

	"github.com/AndreaCostanzo1/http-proxy/http_proxy"
)

func TestLinearRetry_GetNextTimer(t *testing.T) {
	tests := []struct {
		name      string
		options   http_proxy.LinearRetry
		retries   int
		wantTimer int
	}{
		{"Test 0 retries", http_proxy.LinearRetry{Milliseconds: 1000}, 0, 0},
		{"Test 1 retry", http_proxy.LinearRetry{Milliseconds: 1000}, 1, 1000},
		{"Test 5 retries", http_proxy.LinearRetry{Milliseconds: 1000}, 5, 5000},
		{"Test different milliseconds", http_proxy.LinearRetry{Milliseconds: 500}, 5, 2500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.options.GetNextTimer(tt.retries); got != tt.wantTimer {
				t.Errorf("LinearRetry.GetNextTimer() = %v, want %v", got, tt.wantTimer)
			}
		})
	}
}

// ExponentialRetry tests
func TestExponentialRetry_GetNextTimer(t *testing.T) {
	tests := []struct {
		name      string
		options   http_proxy.ExponentialRetry
		retries   int
		wantTimer int
	}{
		{"Test 0 retries", http_proxy.ExponentialRetry{Base: 2, InitialDelayMS: 1000, MaxDelay: 0}, 0, 1000},
		{"Test 1 retry", http_proxy.ExponentialRetry{Base: 2, InitialDelayMS: 1000, MaxDelay: 0}, 1, 2000},
		{"Test 2 retries", http_proxy.ExponentialRetry{Base: 2, InitialDelayMS: 1000, MaxDelay: 0}, 2, 4000},
		{"Test max delay hit", http_proxy.ExponentialRetry{Base: 2, InitialDelayMS: 1000, MaxDelay: 5000}, 3, 5000},
		{"Test no max delay", http_proxy.ExponentialRetry{Base: 2, InitialDelayMS: 1000, MaxDelay: 0}, 3, 8000},
		{"Test different base", http_proxy.ExponentialRetry{Base: 3, InitialDelayMS: 500, MaxDelay: 0}, 2, 4500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.options.GetNextTimer(tt.retries); got != tt.wantTimer {
				t.Errorf("ExponentialRetry.GetNextTimer() = %v, want %v", got, tt.wantTimer)
			}
		})
	}
}
