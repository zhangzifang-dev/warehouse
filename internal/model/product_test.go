package model

import (
	"testing"
)

func TestProduct_TableName(t *testing.T) {
	p := &Product{}
	if p.TableName() != "products" {
		t.Errorf("expected table name 'products', got '%s'", p.TableName())
	}
}

func TestProduct_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		expected bool
	}{
		{"active status", ProductStatusActive, true},
		{"inactive status", ProductStatusInactive, false},
		{"unknown status", 999, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Product{Status: tt.status}
			if p.IsActive() != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, p.IsActive())
			}
		})
	}
}
