package http_proxy_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	http_proxy "github.com/AndreaCostanzo1/http-proxy/pkg"
)

func TestWithGenericInterceptor(t *testing.T) {
	t.Run("WithGenericInterceptor adds generic interceptors", func(t *testing.T) {
		interceptorCalled := false
		interceptor := func(body map[string]interface{}, response *http.Response) error {
			interceptorCalled = true
			return nil
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		}))
		defer server.Close()

		req := http_proxy.NewRequest("GET", server.URL)
		req.WithGenericInterceptor(interceptor)
		resp, err := req.Send()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, got %d", resp.StatusCode)
		}
		if !interceptorCalled {
			t.Errorf("expected generic interceptor to be called")
		}
	})
}

func TestWithStatusCodeInterceptor(t *testing.T) {
	t.Run("WithStatusCodeInterceptor adds status code specific interceptors", func(t *testing.T) {
		interceptorCalled := false
		interceptor := func(body map[string]interface{}, response *http.Response) error {
			interceptorCalled = true
			return nil
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"status":"not found"}`))
		}))
		defer server.Close()

		req := http_proxy.NewRequest("GET", server.URL)
		req.WithStatusCodeInterceptor(http.StatusNotFound, interceptor)
		resp, err := req.Send()

		if err != nil && err.Error() != "404 Not Found" {
			t.Errorf("expected error '404 Not Found', got %v", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status code 404, got %d", resp.StatusCode)
		}
		if !interceptorCalled {
			t.Errorf("expected status code interceptor to be called")
		}
	})
}
