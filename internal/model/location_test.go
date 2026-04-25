package model

import (
	"testing"
)

func TestLocation_TableName(t *testing.T) {
	loc := &Location{}
	if loc.TableName() != "locations" {
		t.Errorf("expected table name 'locations', got '%s'", loc.TableName())
	}
}

func TestLocation_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		expected bool
	}{
		{"active status", LocationStatusActive, true},
		{"inactive status", LocationStatusInactive, false},
		{"other status", 99, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := &Location{Status: tt.status}
			if loc.IsActive() != tt.expected {
				t.Errorf("IsActive() = %v, expected %v", loc.IsActive(), tt.expected)
			}
		})
	}
}

func TestLocation_GenerateCode(t *testing.T) {
	loc := &Location{
		Zone:     "A",
		Shelf:    "01",
		Level:    "02",
		Position: "03",
	}
	expected := "A-01-02-03"
	if loc.GenerateCode() != expected {
		t.Errorf("GenerateCode() = %s, expected %s", loc.GenerateCode(), expected)
	}
}
