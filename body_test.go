package http_proxy_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	http_proxy "github.com/AndreaCostanzo1/http-proxy"
)

func TestSetBody(t *testing.T) {
	t.Run("SetBody with io.Reader", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			if string(body) != "test body" {
				t.Errorf("expected body to be 'test body', got '%s'", body)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		req := http_proxy.NewRequest("POST", server.URL)
		body := bytes.NewBufferString("test body")

		req.SetBody(body)
		resp, err := req.Send()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, got %d", resp.StatusCode)
		}
	})

	t.Run("SetBody with empty io.Reader", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			if string(body) != "" {
				t.Errorf("expected empty body, got '%s'", body)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		req := http_proxy.NewRequest("POST", server.URL)
		body := bytes.NewBufferString("")

		req.SetBody(body)
		resp, err := req.Send()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, got %d", resp.StatusCode)
		}
	})
}

func TestSetJSONBody(t *testing.T) {
	t.Run("SetJSONBody with valid JSON", func(t *testing.T) {
		bodyKey := "key"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var data map[string]string
			json.NewDecoder(r.Body).Decode(&data)
			if data[bodyKey] != "value" {
				t.Errorf("expected JSON key to be 'value', got '%s'", data[bodyKey])
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		req := http_proxy.NewRequest("POST", server.URL)
		data := map[string]string{bodyKey: "value"}

		req.SetJSONBody(data)
		resp, err := req.Send()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, got %d", resp.StatusCode)
		}
	})

	t.Run("SetJSONBody with nil body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var data interface{}
			json.NewDecoder(r.Body).Decode(&data)
			if data != nil {
				t.Errorf("expected nil JSON body, got '%v'", data)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		req := http_proxy.NewRequest("POST", server.URL)

		req.SetJSONBody(nil)
		resp, err := req.Send()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, got %d", resp.StatusCode)
		}
	})

	t.Run("SetJSONBody with invalid JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		req := http_proxy.NewRequest("POST", server.URL)
		data := map[string]interface{}{"key": make(chan int)}

		req.SetJSONBody(data)
		resp, err := req.Send()
		if err == nil {
			t.Errorf("expected an error due to marshal error, got none")
		}
		if resp != nil {
			t.Errorf("expected no response due to marshal error, got %v", resp)
		}
	})
}
