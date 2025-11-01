package httpclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPut(t *testing.T) {
	type testData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected method PUT, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		body, _ := io.ReadAll(r.Body)
		expectedBody := `{"name":"test","value":42}`
		if string(body) != expectedBody {
			t.Errorf("Expected body %s, got %s", expectedBody, string(body))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client := New(server.URL)
	ctx := context.Background()

	data := testData{Name: "test", Value: 42}
	resp, err := client.Put(ctx, "/test", data)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestPutRaw(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected method PUT, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "text/plain" {
			t.Errorf("Expected Content-Type text/plain, got %s", r.Header.Get("Content-Type"))
		}

		body, _ := io.ReadAll(r.Body)
		expectedBody := "raw body content"
		if string(body) != expectedBody {
			t.Errorf("Expected body %s, got %s", expectedBody, string(body))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client := New(server.URL)
	ctx := context.Background()

	resp, err := client.PutRaw(ctx, "/test", "raw body content", "text/plain")
	if err != nil {
		t.Fatalf("PutRaw failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestPutWithNilBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected method PUT, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client := New(server.URL)
	ctx := context.Background()

	resp, err := client.Put(ctx, "/test", nil)
	if err != nil {
		t.Fatalf("Put with nil body failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestDelete(t *testing.T) {
	type deleteData struct {
		Reason string `json:"reason"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected method DELETE, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Deleted"))
	}))
	defer server.Close()

	client := New(server.URL)
	ctx := context.Background()

	data := deleteData{Reason: "testing"}
	resp, err := client.Delete(ctx, "/resource/123", data)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestDeleteWithNilBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected method DELETE, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := New(server.URL)
	ctx := context.Background()

	resp, err := client.Delete(ctx, "/resource/123", nil)
	if err != nil {
		t.Fatalf("Delete with nil body failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}

func TestPutError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	client := New(server.URL)
	ctx := context.Background()

	_, err := client.Put(ctx, "/test", map[string]string{"key": "value"})
	if err == nil {
		t.Fatal("Expected error for 400 status, got nil")
	}

	if !strings.Contains(err.Error(), "Bad Request") {
		t.Errorf("Expected error message to contain 'Bad Request', got: %v", err)
	}
}

func TestDeleteError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	client := New(server.URL)
	ctx := context.Background()

	_, err := client.Delete(ctx, "/resource/999", nil)
	if err == nil {
		t.Fatal("Expected error for 404 status, got nil")
	}

	if !strings.Contains(err.Error(), "Not Found") {
		t.Errorf("Expected error message to contain 'Not Found', got: %v", err)
	}
}

func TestPutWithCustomHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "custom-value" {
			t.Errorf("Expected X-Custom-Header to be 'custom-value', got %s", r.Header.Get("X-Custom-Header"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client := New(server.URL, WithHeaders(map[string]string{
		"X-Custom-Header": "custom-value",
	}))
	ctx := context.Background()

	resp, err := client.Put(ctx, "/test", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("Put with custom headers failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestDeleteWithCustomHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token123" {
			t.Errorf("Expected Authorization header, got %s", r.Header.Get("Authorization"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client := New(server.URL, WithHeaders(map[string]string{
		"Authorization": "Bearer token123",
	}))
	ctx := context.Background()

	resp, err := client.Delete(ctx, "/resource/123", nil)
	if err != nil {
		t.Fatalf("Delete with custom headers failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
