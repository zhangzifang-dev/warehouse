package repository

import (
	"testing"

	"warehouse/internal/config"
)

func TestBuildDSN(t *testing.T) {
	tests := []struct {
		name     string
		dbConfig config.DatabaseConfig
		expected string
	}{
		{
			name: "basic config",
			dbConfig: config.DatabaseConfig{
				Host:     "localhost",
				Port:     3306,
				User:     "root",
				Password: "password",
				Name:     "testdb",
			},
			expected: "root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=true",
		},
		{
			name: "empty password",
			dbConfig: config.DatabaseConfig{
				Host:     "127.0.0.1",
				Port:     3307,
				User:     "admin",
				Password: "",
				Name:     "mydb",
			},
			expected: "admin:@tcp(127.0.0.1:3307)/mydb?charset=utf8mb4&parseTime=true",
		},
		{
			name: "special chars in password",
			dbConfig: config.DatabaseConfig{
				Host:     "db.example.com",
				Port:     3306,
				User:     "user",
				Password: "p@ss!word",
				Name:     "production",
			},
			expected: "user:p@ss!word@tcp(db.example.com:3306)/production?charset=utf8mb4&parseTime=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildDSN(tt.dbConfig)
			if result != tt.expected {
				t.Errorf("BuildDSN() = %s, want %s", result, tt.expected)
			}
		})
	}
}
