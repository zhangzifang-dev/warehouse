package model

import "testing"

func TestSupplier_TableName(t *testing.T) {
	s := &Supplier{}
	if s.TableName() != "suppliers" {
		t.Errorf("expected table name 'suppliers', got '%s'", s.TableName())
	}
}

func TestSupplier_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		expected bool
	}{
		{"active status", SupplierStatusActive, true},
		{"inactive status", SupplierStatusInactive, false},
		{"unknown status", 99, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Supplier{Status: tt.status}
			if s.IsActive() != tt.expected {
				t.Errorf("expected IsActive() to be %v, got %v", tt.expected, s.IsActive())
			}
		})
	}
}

func TestCustomer_TableName(t *testing.T) {
	c := &Customer{}
	if c.TableName() != "customers" {
		t.Errorf("expected table name 'customers', got '%s'", c.TableName())
	}
}

func TestCustomer_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		expected bool
	}{
		{"active status", CustomerStatusActive, true},
		{"inactive status", CustomerStatusInactive, false},
		{"unknown status", 99, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Customer{Status: tt.status}
			if c.IsActive() != tt.expected {
				t.Errorf("expected IsActive() to be %v, got %v", tt.expected, c.IsActive())
			}
		})
	}
}
