package pagination

import "testing"

func TestSetDefault(t *testing.T) {
	tests := []struct {
		name     string
		p        Pagination
		expected Pagination
	}{
		{
			name:     "both zero values",
			p:        Pagination{Page: 0, PageSize: 0},
			expected: Pagination{Page: 1, PageSize: 20},
		},
		{
			name:     "page zero, pageSize set",
			p:        Pagination{Page: 0, PageSize: 10},
			expected: Pagination{Page: 1, PageSize: 10},
		},
		{
			name:     "page set, pageSize zero",
			p:        Pagination{Page: 3, PageSize: 0},
			expected: Pagination{Page: 3, PageSize: 20},
		},
		{
			name:     "both set, no change",
			p:        Pagination{Page: 2, PageSize: 15},
			expected: Pagination{Page: 2, PageSize: 15},
		},
		{
			name:     "negative page",
			p:        Pagination{Page: -1, PageSize: 10},
			expected: Pagination{Page: -1, PageSize: 10},
		},
		{
			name:     "negative pageSize",
			p:        Pagination{Page: 1, PageSize: -5},
			expected: Pagination{Page: 1, PageSize: -5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.p.SetDefault()
			if result.Page != tt.expected.Page {
				t.Errorf("SetDefault() Page = %d, want %d", result.Page, tt.expected.Page)
			}
			if result.PageSize != tt.expected.PageSize {
				t.Errorf("SetDefault() PageSize = %d, want %d", result.PageSize, tt.expected.PageSize)
			}
		})
	}
}

func TestLimit(t *testing.T) {
	tests := []struct {
		name     string
		p        Pagination
		expected int
	}{
		{
			name:     "pageSize 10",
			p:        Pagination{Page: 1, PageSize: 10},
			expected: 10,
		},
		{
			name:     "pageSize 20",
			p:        Pagination{Page: 2, PageSize: 20},
			expected: 20,
		},
		{
			name:     "pageSize 50",
			p:        Pagination{Page: 3, PageSize: 50},
			expected: 50,
		},
		{
			name:     "pageSize 1",
			p:        Pagination{Page: 1, PageSize: 1},
			expected: 1,
		},
		{
			name:     "pageSize 0",
			p:        Pagination{Page: 1, PageSize: 0},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.p.Limit()
			if result != tt.expected {
				t.Errorf("Limit() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestOffset(t *testing.T) {
	tests := []struct {
		name     string
		p        Pagination
		expected int
	}{
		{
			name:     "page 1, pageSize 10",
			p:        Pagination{Page: 1, PageSize: 10},
			expected: 0,
		},
		{
			name:     "page 2, pageSize 10",
			p:        Pagination{Page: 2, PageSize: 10},
			expected: 10,
		},
		{
			name:     "page 3, pageSize 20",
			p:        Pagination{Page: 3, PageSize: 20},
			expected: 40,
		},
		{
			name:     "page 5, pageSize 15",
			p:        Pagination{Page: 5, PageSize: 15},
			expected: 60,
		},
		{
			name:     "page 1, pageSize 1",
			p:        Pagination{Page: 1, PageSize: 1},
			expected: 0,
		},
		{
			name:     "page 0, pageSize 10",
			p:        Pagination{Page: 0, PageSize: 10},
			expected: -10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.p.Offset()
			if result != tt.expected {
				t.Errorf("Offset() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestSetTotal(t *testing.T) {
	tests := []struct {
		name          string
		p             Pagination
		totalData     int
		expectedTotal int
		expectedPages int
	}{
		{
			name:          "total data less than page size",
			p:             Pagination{Page: 1, PageSize: 20},
			totalData:     15,
			expectedTotal: 15,
			expectedPages: 1,
		},
		{
			name:          "total data equals page size",
			p:             Pagination{Page: 1, PageSize: 20},
			totalData:     20,
			expectedTotal: 20,
			expectedPages: 1,
		},
		{
			name:          "total data greater than page size",
			p:             Pagination{Page: 1, PageSize: 20},
			totalData:     45,
			expectedTotal: 45,
			expectedPages: 3,
		},
		{
			name:          "total data requires rounding up",
			p:             Pagination{Page: 1, PageSize: 20},
			totalData:     41,
			expectedTotal: 41,
			expectedPages: 3,
		},
		{
			name:          "exact multiple of page size",
			p:             Pagination{Page: 1, PageSize: 20},
			totalData:     100,
			expectedTotal: 100,
			expectedPages: 5,
		},
		{
			name:          "zero total data",
			p:             Pagination{Page: 1, PageSize: 20},
			totalData:     0,
			expectedTotal: 0,
			expectedPages: 0,
		},
		{
			name:          "pageSize 1, total data 5",
			p:             Pagination{Page: 1, PageSize: 1},
			totalData:     5,
			expectedTotal: 5,
			expectedPages: 5,
		},
		{
			name:          "large page size",
			p:             Pagination{Page: 1, PageSize: 100},
			totalData:     250,
			expectedTotal: 250,
			expectedPages: 3,
		},
		{
			name:          "total data one less than multiple",
			p:             Pagination{Page: 1, PageSize: 20},
			totalData:     59,
			expectedTotal: 59,
			expectedPages: 3,
		},
		{
			name:          "total data one more than multiple",
			p:             Pagination{Page: 1, PageSize: 20},
			totalData:     61,
			expectedTotal: 61,
			expectedPages: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.p.SetTotal(tt.totalData)
			if result.TotalData != tt.expectedTotal {
				t.Errorf("SetTotal() TotalData = %d, want %d", result.TotalData, tt.expectedTotal)
			}
			if result.TotalPage != tt.expectedPages {
				t.Errorf("SetTotal() TotalPage = %d, want %d", result.TotalPage, tt.expectedPages)
			}
			// Verify the original struct is also modified
			if tt.p.TotalData != tt.expectedTotal {
				t.Errorf("SetTotal() did not modify original TotalData, got %d, want %d", tt.p.TotalData, tt.expectedTotal)
			}
			if tt.p.TotalPage != tt.expectedPages {
				t.Errorf("SetTotal() did not modify original TotalPage, got %d, want %d", tt.p.TotalPage, tt.expectedPages)
			}
		})
	}
}

func TestPaginationIntegration(t *testing.T) {
	// Test a complete workflow
	p := Pagination{Page: 0, PageSize: 0}

	// Step 1: Set defaults
	result := p.SetDefault()
	if result.Page != 1 {
		t.Errorf("Integration: SetDefault() Page = %d, want 1", result.Page)
	}
	if result.PageSize != 20 {
		t.Errorf("Integration: SetDefault() PageSize = %d, want 20", result.PageSize)
	}

	// Step 2: Calculate offset for page 2
	p.Page = 2
	offset := p.Offset()
	if offset != 20 {
		t.Errorf("Integration: Offset() = %d, want 20", offset)
	}

	// Step 3: Set total and verify total pages
	p.SetTotal(45)
	if p.TotalData != 45 {
		t.Errorf("Integration: TotalData = %d, want 45", p.TotalData)
	}
	if p.TotalPage != 3 {
		t.Errorf("Integration: TotalPage = %d, want 3", p.TotalPage)
	}

	// Step 4: Verify limit
	limit := p.Limit()
	if limit != 20 {
		t.Errorf("Integration: Limit() = %d, want 20", limit)
	}
}
