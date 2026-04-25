package model

import (
	"testing"
)

func TestCategory_TableName(t *testing.T) {
	cat := &Category{}
	if cat.TableName() != "categories" {
		t.Errorf("expected table name 'categories', got '%s'", cat.TableName())
	}
}

func TestCategory_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		expected bool
	}{
		{"active status", CategoryStatusActive, true},
		{"inactive status", CategoryStatusInactive, false},
		{"other status", 99, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cat := &Category{Status: tt.status}
			if cat.IsActive() != tt.expected {
				t.Errorf("IsActive() = %v, expected %v", cat.IsActive(), tt.expected)
			}
		})
	}
}
