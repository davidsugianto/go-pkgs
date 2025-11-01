package main

import (
	"encoding/json"
	"fmt"

	"github.com/davidsugianto/go-pkgs/pagination"
)

// Simulate database records
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Simulate a database query function
func fetchUsers(offset, limit int) []User {
	// In real implementation, this would query the database
	// SELECT * FROM users LIMIT ? OFFSET ?
	users := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}
	if offset >= len(users) {
		return []User{}
	}
	if offset+limit > len(users) {
		return users[offset:]
	}
	return users[offset : offset+limit]
}

// Simulate counting total records
func countUsers() int {
	// In real implementation: SELECT COUNT(*) FROM users
	return 150 // Simulated total
}

func main() {
	fmt.Println("=== Pagination Package Example ===")
	fmt.Println()

	// Example 1: Basic pagination with defaults
	fmt.Println("1. Basic Pagination with Defaults")
	p := pagination.Pagination{
		Page:     0, // Will be set to 1 by SetDefault()
		PageSize: 0, // Will be set to 20 by SetDefault()
	}
	p = p.SetDefault()
	fmt.Printf("   Page: %d, PageSize: %d\n", p.Page, p.PageSize)
	fmt.Printf("   Limit: %d, Offset: %d\n", p.Limit(), p.Offset())
	fmt.Println()

	// Example 2: Custom pagination
	fmt.Println("2. Custom Pagination")
	p2 := pagination.Pagination{
		Page:     2,
		PageSize: 10,
	}
	fmt.Printf("   Page: %d, PageSize: %d\n", p2.Page, p2.PageSize)
	fmt.Printf("   Limit: %d, Offset: %d\n", p2.Limit(), p2.Offset())
	fmt.Println()

	// Example 3: Complete workflow - query with pagination
	fmt.Println("3. Complete Workflow - Query with Pagination")
	p3 := pagination.Pagination{
		Page:     3,
		PageSize: 25,
	}

	// Calculate offset and limit for database query
	offset := p3.Offset()
	limit := p3.Limit()
	fmt.Printf("   Query parameters: OFFSET %d LIMIT %d\n", offset, limit)

	// Simulate database query
	users := fetchUsers(offset, limit)
	fmt.Printf("   Fetched %d users\n", len(users))

	// Set total after querying count
	totalUsers := countUsers()
	p3 = p3.SetTotal(totalUsers)
	fmt.Printf("   Total Data: %d\n", p3.TotalData)
	fmt.Printf("   Total Pages: %d\n", p3.TotalPage)
	fmt.Println()

	// Example 4: Pagination response for API
	fmt.Println("4. Pagination Response for API")
	p4 := pagination.Pagination{
		Page:     1,
		PageSize: 20,
	}
	p4 = p4.SetTotal(145) // 145 total records

	// Create a response object
	response := map[string]interface{}{
		"pagination": p4,
		"data":       []User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}},
	}

	jsonResponse, _ := json.MarshalIndent(response, "   ", "  ")
	fmt.Println("   JSON Response:")
	fmt.Println(string(jsonResponse))
	fmt.Println()

	// Example 5: Different page sizes
	fmt.Println("5. Different Page Sizes")
	pages := []int{1, 2, 3}
	pageSize := 15
	total := 47

	for _, page := range pages {
		p := pagination.Pagination{
			Page:     page,
			PageSize: pageSize,
		}
		p = p.SetTotal(total)
		fmt.Printf("   Page %d: Offset=%d, Limit=%d, TotalPages=%d\n",
			p.Page, p.Offset(), p.Limit(), p.TotalPage)
	}
	fmt.Println()

	// Example 6: Edge case - zero or empty results
	fmt.Println("6. Edge Case - Zero Results")
	p6 := pagination.Pagination{
		Page:     1,
		PageSize: 20,
	}
	p6 = p6.SetTotal(0)
	fmt.Printf("   Total Data: %d\n", p6.TotalData)
	fmt.Printf("   Total Pages: %d\n", p6.TotalPage)
	fmt.Printf("   Offset: %d, Limit: %d\n", p6.Offset(), p6.Limit())
	fmt.Println()

	// Example 7: Large dataset
	fmt.Println("7. Large Dataset Example")
	p7 := pagination.Pagination{
		Page:     50,
		PageSize: 100,
	}
	p7 = p7.SetTotal(10000)
	fmt.Printf("   Page: %d of %d\n", p7.Page, p7.TotalPage)
	fmt.Printf("   Offset: %d, Limit: %d\n", p7.Offset(), p7.Limit())
	fmt.Printf("   Showing records %d to %d of %d\n",
		p7.Offset()+1, p7.Offset()+p7.Limit(), p7.TotalData)
	fmt.Println()

	fmt.Println("=== Example Complete ===")
}
