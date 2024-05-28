package http_proxy_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	http_proxy "github.com/AndreaCostanzo1/http-proxy"
)

func TestWithContext(t *testing.T) {
	t.Run("WithContext sets context correctly", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		req := http_proxy.NewRequest("GET", server.URL)
		ctx := context.WithValue(context.Background(), "key", "value")

		req.WithContext(ctx)
		resp, err := req.Send()
		if resp.Request.Context().Value("key") != "value" {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, got %d", resp.StatusCode)
		}
	})

	t.Run("WithContext with canceled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		req := http_proxy.NewRequest("GET", server.URL)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		req.WithContext(ctx)
		resp, err := req.Send()

		if err == nil {
			t.Errorf("expected an error due to canceled context, got none")
		}
		if resp != nil {
			t.Errorf("expected no response due to canceled context, got %v", resp)
		}
	})
}
