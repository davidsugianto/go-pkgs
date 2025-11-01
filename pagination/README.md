# Pagination - Pagination Helper

A lightweight pagination package for Go that helps you handle pagination logic in your applications.

## Installation

```bash
go get github.com/davidsugianto/go-pkgs/pagination
```

## Quick Start

```go
package main

import (
    "github.com/davidsugianto/go-pkgs/pagination"
)

func main() {
    // Create pagination with defaults (Page=1, PageSize=20)
    p := pagination.Pagination{}
    p = p.SetDefault()
    
    // Use for database queries
    offset := p.Offset()  // 0
    limit := p.Limit()    // 20
    
    // After querying total records
    p = p.SetTotal(145)
    // p.TotalPage will be calculated automatically: 8
}
```

## Features

- ✅ **Default Values** - Automatically sets Page=1 and PageSize=20 when zero
- ✅ **Offset/Limit Calculation** - Easy calculation for database queries
- ✅ **Total Pages** - Automatically calculates total pages from total data
- ✅ **JSON Support** - JSON tags included for API responses

## Usage

### Basic Usage

```go
// Custom pagination
p := pagination.Pagination{
    Page:     2,
    PageSize: 10,
}

offset := p.Offset()  // 10
limit := p.Limit()    // 10
```

### With Defaults

```go
// Zero values will be set to defaults
p := pagination.Pagination{
    Page:     0, // Will become 1
    PageSize: 0, // Will become 20
}
p = p.SetDefault()
// p.Page = 1, p.PageSize = 20
```

### Database Query Example

```go
// For SQL queries: SELECT * FROM users LIMIT ? OFFSET ?
p := pagination.Pagination{
    Page:     3,
    PageSize: 25,
}

offset := p.Offset()  // 50
limit := p.Limit()    // 25

// Query database
users := db.Query("SELECT * FROM users LIMIT ? OFFSET ?", limit, offset)

// After getting total count
totalUsers := db.Count("SELECT COUNT(*) FROM users")
p = p.SetTotal(totalUsers)
// p.TotalPage will be calculated automatically
```

### API Response Example

```go
p := pagination.Pagination{
    Page:     1,
    PageSize: 20,
}
p = p.SetTotal(145)

// The Pagination struct has JSON tags and can be used directly in responses
response := map[string]interface{}{
    "pagination": p,
    "data":       yourData,
}
// Returns: {"pagination":{"page":1,"page_size":20,"total_data":145,"total_page":8},...}
```

## API Reference

### Methods

#### `SetDefault() Pagination`

Sets default values when Page or PageSize is zero:
- Page: 0 → 1
- PageSize: 0 → 20

#### `Limit() int`

Returns the page size (limit for queries).

#### `Offset() int`

Calculates and returns the offset: `(Page - 1) * PageSize`

#### `SetTotal(totalData int) Pagination`

Sets the total number of records and calculates total pages:
- `TotalData`: Set to `totalData`
- `TotalPage`: Calculated as `(totalData + PageSize - 1) / PageSize`

### Struct Fields

```go
type Pagination struct {
    Page      int `form:"page" json:"page"`           // Current page (1-indexed)
    PageSize  int `form:"page_size" json:"page_size"` // Items per page
    TotalData int `json:"total_data"`                 // Total number of records
    TotalPage int `json:"total_page"`                 // Total number of pages
}
```

## Examples

See the `example/` directory for a complete working example.

```bash
cd example
go run main.go
```

