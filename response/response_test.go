package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSON(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		data       interface{}
		wantCode   int
		wantData   string
	}{
		{
			name:       "success with map",
			statusCode: http.StatusOK,
			data:       map[string]string{"message": "success"},
			wantCode:   http.StatusOK,
			wantData:   `"message":"success"`,
		},
		{
			name:       "success with struct",
			statusCode: http.StatusOK,
			data: struct {
				ID int `json:"id"`
			}{ID: 123},
			wantCode: http.StatusOK,
			wantData: `"id":123`,
		},
		{
			name:       "created status",
			statusCode: http.StatusCreated,
			data:       map[string]int{"id": 456},
			wantCode:   http.StatusCreated,
			wantData:   `"id":456`,
		},
		{
			name:       "nil data",
			statusCode: http.StatusOK,
			data:       nil,
			wantCode:   http.StatusOK,
			wantData:   `"data":null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := JSON(w, tt.statusCode, tt.data)
			if err != nil {
				t.Errorf("JSON() error = %v", err)
				return
			}

			if w.Code != tt.wantCode {
				t.Errorf("JSON() statusCode = %v, want %v", w.Code, tt.wantCode)
			}

			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("JSON() Content-Type = %v, want application/json", w.Header().Get("Content-Type"))
			}

			body := w.Body.String()
			body = strings.TrimSpace(body) // Remove trailing newline from json.Encode

			// Verify JSON structure
			var resp Response
			if err := json.Unmarshal([]byte(body), &resp); err != nil {
				t.Errorf("JSON() invalid JSON response: %v, body: %v", err, body)
				return
			}

			if resp.Code != tt.statusCode {
				t.Errorf("JSON() code = %v, want %v", resp.Code, tt.statusCode)
			}

			// For nil data test, check that data field is omitted (due to omitempty)
			if tt.name == "nil data" {
				if strings.Contains(body, `"data"`) {
					t.Errorf("JSON() body should NOT contain 'data' field when data is nil (omitempty), got: %v", body)
				}
			} else if tt.data != nil {
				// Verify data is present in JSON
				if !strings.Contains(body, tt.wantData) {
					t.Errorf("JSON() body = %v, want to contain %v", body, tt.wantData)
				}
			}
		})
	}
}

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"message": "success"}

	err := Success(w, data)
	if err != nil {
		t.Errorf("Success() error = %v", err)
		return
	}

	if w.Code != http.StatusOK {
		t.Errorf("Success() statusCode = %v, want %v", w.Code, http.StatusOK)
	}

	body := strings.TrimSpace(w.Body.String())

	var resp Response
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Errorf("Success() invalid JSON response: %v", err)
		return
	}

	if resp.Code != http.StatusOK {
		t.Errorf("Success() code = %v, want %v", resp.Code, http.StatusOK)
	}

	if !strings.Contains(body, `"message"`) || !strings.Contains(body, `"success"`) {
		t.Errorf("Success() body should contain message data, got: %v", body)
	}
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]int{"id": 123}

	err := Created(w, data)
	if err != nil {
		t.Errorf("Created() error = %v", err)
		return
	}

	if w.Code != http.StatusCreated {
		t.Errorf("Created() statusCode = %v, want %v", w.Code, http.StatusCreated)
	}

	body := strings.TrimSpace(w.Body.String())

	var resp Response
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Errorf("Created() invalid JSON response: %v", err)
		return
	}

	if resp.Code != http.StatusCreated {
		t.Errorf("Created() code = %v, want %v", resp.Code, http.StatusCreated)
	}
}

func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	NoContent(w)

	if w.Code != http.StatusNoContent {
		t.Errorf("NoContent() statusCode = %v, want %v", w.Code, http.StatusNoContent)
	}

	if w.Body.Len() != 0 {
		t.Errorf("NoContent() body should be empty, got length %v", w.Body.Len())
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		err        error
		wantCode   int
		wantError  string
	}{
		{
			name:       "with error",
			statusCode: http.StatusBadRequest,
			err:        errors.New("validation failed"),
			wantCode:   http.StatusBadRequest,
			wantError:  "validation failed",
		},
		{
			name:       "nil error",
			statusCode: http.StatusInternalServerError,
			err:        nil,
			wantCode:   http.StatusInternalServerError,
			wantError:  "",
		},
		{
			name:       "not found",
			statusCode: http.StatusNotFound,
			err:        errors.New("resource not found"),
			wantCode:   http.StatusNotFound,
			wantError:  "resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := Error(w, tt.statusCode, tt.err)
			if err != nil {
				t.Errorf("Error() error = %v", err)
				return
			}

			if w.Code != tt.wantCode {
				t.Errorf("Error() statusCode = %v, want %v", w.Code, tt.wantCode)
			}

			body := strings.TrimSpace(w.Body.String())

			var resp Response
			if err := json.Unmarshal([]byte(body), &resp); err != nil {
				t.Errorf("Error() invalid JSON response: %v, body: %v", err, body)
				return
			}

			if resp.Code != tt.statusCode {
				t.Errorf("Error() code = %v, want %v", resp.Code, tt.statusCode)
			}

			if tt.wantError != "" && resp.Error != tt.wantError {
				t.Errorf("Error() error message = %v, want %v", resp.Error, tt.wantError)
			}
		})
	}
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	err := BadRequest(w, errors.New("invalid input"))
	if err != nil {
		t.Errorf("BadRequest() error = %v", err)
		return
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("BadRequest() statusCode = %v, want %v", w.Code, http.StatusBadRequest)
	}

	body := strings.TrimSpace(w.Body.String())

	var resp Response
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Errorf("BadRequest() invalid JSON response: %v", err)
		return
	}

	if resp.Code != http.StatusBadRequest {
		t.Errorf("BadRequest() code = %v, want %v", resp.Code, http.StatusBadRequest)
	}

	if !strings.Contains(body, "invalid input") {
		t.Errorf("BadRequest() body should contain error message, got: %v", body)
	}
}

func TestUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	err := Unauthorized(w, errors.New("authentication required"))
	if err != nil {
		t.Errorf("Unauthorized() error = %v", err)
		return
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Unauthorized() statusCode = %v, want %v", w.Code, http.StatusUnauthorized)
	}

	body := strings.TrimSpace(w.Body.String())

	var resp Response
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Errorf("Unauthorized() invalid JSON response: %v", err)
		return
	}

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Unauthorized() code = %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

func TestForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	err := Forbidden(w, errors.New("access denied"))
	if err != nil {
		t.Errorf("Forbidden() error = %v", err)
		return
	}

	if w.Code != http.StatusForbidden {
		t.Errorf("Forbidden() statusCode = %v, want %v", w.Code, http.StatusForbidden)
	}

	body := strings.TrimSpace(w.Body.String())

	var resp Response
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Errorf("Forbidden() invalid JSON response: %v", err)
		return
	}

	if resp.Code != http.StatusForbidden {
		t.Errorf("Forbidden() code = %v, want %v", resp.Code, http.StatusForbidden)
	}
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	err := NotFound(w, errors.New("resource not found"))
	if err != nil {
		t.Errorf("NotFound() error = %v", err)
		return
	}

	if w.Code != http.StatusNotFound {
		t.Errorf("NotFound() statusCode = %v, want %v", w.Code, http.StatusNotFound)
	}

	body := strings.TrimSpace(w.Body.String())

	var resp Response
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Errorf("NotFound() invalid JSON response: %v", err)
		return
	}

	if resp.Code != http.StatusNotFound {
		t.Errorf("NotFound() code = %v, want %v", resp.Code, http.StatusNotFound)
	}
}

func TestInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	err := InternalServerError(w, errors.New("internal error"))
	if err != nil {
		t.Errorf("InternalServerError() error = %v", err)
		return
	}

	if w.Code != http.StatusInternalServerError {
		t.Errorf("InternalServerError() statusCode = %v, want %v", w.Code, http.StatusInternalServerError)
	}

	body := strings.TrimSpace(w.Body.String())

	var resp Response
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Errorf("InternalServerError() invalid JSON response: %v", err)
		return
	}

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("InternalServerError() code = %v, want %v", resp.Code, http.StatusInternalServerError)
	}
}

func TestStatusCode(t *testing.T) {
	w := httptest.NewRecorder()
	message := "custom error message"
	err := StatusCode(w, http.StatusConflict, message)
	if err != nil {
		t.Errorf("StatusCode() error = %v", err)
		return
	}

	if w.Code != http.StatusConflict {
		t.Errorf("StatusCode() statusCode = %v, want %v", w.Code, http.StatusConflict)
	}

	body := strings.TrimSpace(w.Body.String())

	var resp Response
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Errorf("StatusCode() invalid JSON response: %v", err)
		return
	}

	if resp.Code != http.StatusConflict {
		t.Errorf("StatusCode() code = %v, want %v", resp.Code, http.StatusConflict)
	}

	if resp.Error != message {
		t.Errorf("StatusCode() error message = %v, want %v", resp.Error, message)
	}
}

func TestResponseStruct(t *testing.T) {
	tests := []struct {
		name    string
		resp    Response
		wantHas string
	}{
		{
			name:    "with data",
			resp:    Response{Code: 200, Data: map[string]string{"key": "value"}},
			wantHas: `"data"`,
		},
		{
			name:    "with error",
			resp:    Response{Code: 400, Error: "error message"},
			wantHas: `"error"`,
		},
		{
			name:    "both data and error",
			resp:    Response{Code: 200, Data: "data", Error: "error"},
			wantHas: `"code"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(tt.resp.Code)

			err := json.NewEncoder(w).Encode(tt.resp)
			if err != nil {
				t.Errorf("Response encoding error = %v", err)
				return
			}

			body := strings.TrimSpace(w.Body.String())

			var resp Response
			if err := json.Unmarshal([]byte(body), &resp); err != nil {
				t.Errorf("Response unmarshaling error = %v, body: %v", err, body)
				return
			}

			if resp.Code != tt.resp.Code {
				t.Errorf("Response code = %v, want %v", resp.Code, tt.resp.Code)
			}

			if tt.name == "with data" && resp.Data == nil {
				t.Errorf("Response should have data field")
			}

			if tt.name == "with error" && resp.Error == "" {
				t.Errorf("Response should have error field")
			}
		})
	}
}

func TestIntegration(t *testing.T) {
	// Test a complete workflow
	w := httptest.NewRecorder()

	// 1. Success response
	err := Success(w, map[string]string{"status": "ok"})
	if err != nil {
		t.Fatalf("Integration: Success() error = %v", err)
	}
	if w.Code != http.StatusOK {
		t.Errorf("Integration: Expected status 200, got %v", w.Code)
	}

	// 2. Error response
	w2 := httptest.NewRecorder()
	err = BadRequest(w2, errors.New("invalid request"))
	if err != nil {
		t.Fatalf("Integration: BadRequest() error = %v", err)
	}
	if w2.Code != http.StatusBadRequest {
		t.Errorf("Integration: Expected status 400, got %v", w2.Code)
	}

	// 3. No content
	w3 := httptest.NewRecorder()
	NoContent(w3)
	if w3.Code != http.StatusNoContent {
		t.Errorf("Integration: Expected status 204, got %v", w3.Code)
	}
}
