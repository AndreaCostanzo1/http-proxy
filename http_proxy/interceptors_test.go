package http_proxy_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	http_proxy "github.com/AndreaCostanzo1/http-proxy/http_proxy"
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

func TestJsonFormatKeyAndParsedBody(t *testing.T) {
	t.Run("validates JSON FORMAT key and parsed body", func(t *testing.T) {
		expectedStatus := "ok"
		var receivedBody map[string]interface{}

		interceptor := func(body map[string]interface{}, response *http.Response) error {
			receivedBody = body
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
		if receivedBody == nil {
			t.Errorf("expected body to be parsed, got nil")
		} else {
			format, formatExists := receivedBody[http_proxy.FORMAT_TYPE]
			if !formatExists {
				t.Errorf("expected FORMAT key in parsed body")
			}
			if format != http_proxy.FORMAT_JSON {
				t.Errorf("expected FORMAT to be 'JSON', got %v", format)
			}
			status, statusExists := receivedBody["status"]
			if !statusExists {
				t.Errorf("expected 'status' key in parsed body")
			}
			if status != expectedStatus {
				t.Errorf("expected status to be '%s', got '%s'", expectedStatus, status)
			}
		}
	})
}
