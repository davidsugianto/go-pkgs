package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(Response{
		Code: statusCode,
		Data: data,
	})
}

func Success(w http.ResponseWriter, data interface{}) error {
	return JSON(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, data interface{}) error {
	return JSON(w, http.StatusCreated, data)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Error(w http.ResponseWriter, statusCode int, err error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	return json.NewEncoder(w).Encode(Response{
		Code:  statusCode,
		Error: errMsg,
	})
}

func BadRequest(w http.ResponseWriter, err error) error {
	return Error(w, http.StatusBadRequest, err)
}

func Unauthorized(w http.ResponseWriter, err error) error {
	return Error(w, http.StatusUnauthorized, err)
}

func Forbidden(w http.ResponseWriter, err error) error {
	return Error(w, http.StatusForbidden, err)
}

func NotFound(w http.ResponseWriter, err error) error {
	return Error(w, http.StatusNotFound, err)
}

func InternalServerError(w http.ResponseWriter, err error) error {
	return Error(w, http.StatusInternalServerError, err)
}

func StatusCode(w http.ResponseWriter, statusCode int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(Response{
		Code:  statusCode,
		Error: message,
	})
}
