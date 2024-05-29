package http_proxy_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	http_proxy "github.com/AndreaCostanzo1/http-proxy/http_proxy"
)

func TestAddHeader(t *testing.T) {
	t.Run("AddHeader appends header values", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			values := r.Header["X-Custom-Header"]
			expectedValues := []string{"value1", "value2"}
			if len(values) != len(expectedValues) {
				t.Errorf("expected %d values, got %d", len(expectedValues), len(values))
			}
			for i, v := range values {
				if v != expectedValues[i] {
					t.Errorf("expected value '%s', got '%s'", expectedValues[i], v)
				}
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		req := http_proxy.NewRequest("GET", server.URL)
		req.AddHeader("X-Custom-Header", "value1")
		req.AddHeader("X-Custom-Header", "value2")
		resp, err := req.Send()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, got %d", resp.StatusCode)
		}
	})
}

func TestSetHeader(t *testing.T) {
	t.Run("SetHeader sets single header value", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value := r.Header.Get("X-Custom-Header")
			if value != "value" {
				t.Errorf("expected header value 'value', got '%s'", value)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		req := http_proxy.NewRequest("GET", server.URL)
		req.SetHeader("X-Custom-Header", "value")
		resp, err := req.Send()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, got %d", resp.StatusCode)
		}
	})
}

func TestSetHeaders(t *testing.T) {
	t.Run("SetHeaders sets multiple headers", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expectedHeaders := map[string]string{
				"Header-One": "value1",
				"Header-Two": "value2",
			}
			for key, expectedValue := range expectedHeaders {
				if value := r.Header.Get(key); value != expectedValue {
					t.Errorf("expected header '%s' to have value '%s', got '%s'", key, expectedValue, value)
				}
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		req := http_proxy.NewRequest("GET", server.URL)
		headers := map[string]string{
			"Header-One": "value1",
			"Header-Two": "value2",
		}
		req.SetHeaders(headers)
		resp, err := req.Send()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, got %d", resp.StatusCode)
		}
	})
}

func TestSetMultiValueHeaders(t *testing.T) {
	t.Run("SetMultiValueHeaders sets headers with multiple values", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expectedHeaders := map[string][]string{
				"Header-One": {"value1", "value2"},
				"Header-Two": {"value3", "value4"},
			}
			for key, expectedValues := range expectedHeaders {
				values := r.Header[key]
				if len(values) != len(expectedValues) {
					t.Errorf("expected %d values for header '%s', got %d", len(expectedValues), key, len(values))
				}
				for i, v := range values {
					if v != expectedValues[i] {
						t.Errorf("expected value '%s' for header '%s', got '%s'", expectedValues[i], key, v)
					}
				}
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		req := http_proxy.NewRequest("GET", server.URL)
		headers := map[string][]string{
			"Header-One": {"value1", "value2"},
			"Header-Two": {"value3", "value4"},
		}
		req.SetMultiValueHeaders(headers)
		resp, err := req.Send()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, got %d", resp.StatusCode)
		}
	})
}

func TestSetJWTAuthToken(t *testing.T) {
	t.Run("SetJWTAuthToken sets the Authorization header", func(t *testing.T) {
		expectedToken := "some-jwt-token"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			expectedAuthHeader := "Bearer " + expectedToken
			if authHeader != expectedAuthHeader {
				t.Errorf("expected Authorization header to be '%s', got '%s'", expectedAuthHeader, authHeader)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		req := http_proxy.NewRequest("GET", server.URL)
		req.SetJWTAuthToken(expectedToken)
		resp, err := req.Send()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, got %d", resp.StatusCode)
		}
	})
}
