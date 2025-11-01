package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/davidsugianto/go-pkgs/response"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	http.HandleFunc("/success", successHandler)
	http.HandleFunc("/created", createdHandler)
	http.HandleFunc("/not-found", notFoundHandler)
	http.HandleFunc("/bad-request", badRequestHandler)
	http.HandleFunc("/internal-error", internalErrorHandler)
	http.HandleFunc("/no-content", noContentHandler)
	http.HandleFunc("/custom", customHandler)

	log.Println("Server starting on :8080")
	log.Println("Try these endpoints:")
	log.Println("  GET http://localhost:8080/success")
	log.Println("  GET http://localhost:8080/created")
	log.Println("  GET http://localhost:8080/not-found")
	log.Println("  GET http://localhost:8080/bad-request")
	log.Println("  GET http://localhost:8080/internal-error")
	log.Println("  GET http://localhost:8080/no-content")
	log.Println("  GET http://localhost:8080/custom")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func successHandler(w http.ResponseWriter, r *http.Request) {
	user := User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	if err := response.Success(w, user); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func createdHandler(w http.ResponseWriter, r *http.Request) {
	newUser := User{
		ID:    2,
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}

	if err := response.Created(w, newUser); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	if err := response.NotFound(w, errors.New("user not found")); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func badRequestHandler(w http.ResponseWriter, r *http.Request) {
	if err := response.BadRequest(w, errors.New("invalid input parameters")); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func internalErrorHandler(w http.ResponseWriter, r *http.Request) {
	if err := response.InternalServerError(w, errors.New("database connection failed")); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func noContentHandler(w http.ResponseWriter, r *http.Request) {
	response.NoContent(w)
}

func customHandler(w http.ResponseWriter, r *http.Request) {
	// Using JSON for custom status code
	if err := response.JSON(w, http.StatusAccepted, map[string]string{
		"message": "Request accepted for processing",
	}); err != nil {
		log.Printf("Error writing response: %v", err)
	}

	// Or using StatusCode for error messages
	// response.StatusCode(w, http.StatusConflict, "resource already exists")
}
