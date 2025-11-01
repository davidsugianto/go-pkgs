package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/davidsugianto/go-pkgs/httpclient"
)

// Example data structures
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Example using JSONPlaceholder API (a free fake REST API for testing)
	baseURL := "https://jsonplaceholder.typicode.com"

	// Create a client with custom timeout and headers
	client := httpclient.New(
		baseURL,
		httpclient.WithTimeout(30*time.Second),
		httpclient.WithHeaders(map[string]string{
			"User-Agent": "go-pkgs-httpclient/1.0",
		}),
	)

	ctx := context.Background()

	fmt.Println("=== HTTP Client Example ===")
	fmt.Println()

	// Example 1: GET request
	fmt.Println("1. GET Request - Fetch user by ID")
	resp, err := client.Get(ctx, "/users/1", nil)
	if err != nil {
		log.Fatalf("GET failed: %v", err)
	}
	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}
	fmt.Printf("   User: %+v\n\n", user)

	// Example 2: POST request
	fmt.Println("2. POST Request - Create new user")
	newUser := CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	resp, err = client.Post(ctx, "/users", newUser)
	if err != nil {
		log.Fatalf("POST failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("   Response: %s\n\n", string(body))

	// Example 3: PUT request
	fmt.Println("3. PUT Request - Update user")
	updateData := UpdateUserRequest{
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}

	resp, err = client.Put(ctx, "/users/1", updateData)
	if err != nil {
		log.Fatalf("PUT failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("   Response: %s\n\n", string(body))

	// Example 4: PUT with raw body
	fmt.Println("4. PUT Request (Raw) - Update with custom content type")
	xmlBody := `<user><name>Bob</name><email>bob@example.com</email></user>`
	resp, err = client.PutRaw(ctx, "/users/1", xmlBody, "application/xml")
	if err != nil {
		log.Fatalf("PUT Raw failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("   Response: %s\n\n", string(body))

	// Example 5: DELETE request
	fmt.Println("5. DELETE Request - Delete user")
	resp, err = client.Delete(ctx, "/users/1", nil)
	if err != nil {
		log.Fatalf("DELETE failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("   Response Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response Body: %s\n\n", string(body))

	// Example 6: DELETE with body
	fmt.Println("6. DELETE Request with body")
	deleteReason := map[string]string{
		"reason": "No longer needed",
	}
	resp, err = client.Delete(ctx, "/posts/1", deleteReason)
	if err != nil {
		log.Fatalf("DELETE with body failed: %v", err)
	}
	defer resp.Body.Close()

	_, _ = io.ReadAll(resp.Body)
	fmt.Printf("   Response Status: %d\n\n", resp.StatusCode)

	// Example 7: Using context with timeout
	fmt.Println("7. GET Request with custom context timeout")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err = client.Get(timeoutCtx, "/users/2", nil)
	if err != nil {
		// Note: This might fail if the request takes longer than 5 seconds
		fmt.Printf("   Error (expected if timeout): %v\n", err)
	} else {
		defer resp.Body.Close()
		var user2 User
		json.NewDecoder(resp.Body).Decode(&user2)
		fmt.Printf("   User: %+v\n", user2)
	}

	fmt.Println("\n=== Example Complete ===")
}
