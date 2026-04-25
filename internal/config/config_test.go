package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() string
		expectError bool
		validate    func(t *testing.T, cfg *Config)
	}{
		{
			name: "load valid config",
			setup: func() string {
				return filepath.Join("..", "..", "config", "config.yaml")
			},
			expectError: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Server.Port != 8080 {
					t.Errorf("expected server port 8080, got %d", cfg.Server.Port)
				}
				if cfg.Server.Mode != "debug" {
					t.Errorf("expected server mode debug, got %s", cfg.Server.Mode)
				}
				if cfg.Database.Driver != "mysql" {
					t.Errorf("expected database driver mysql, got %s", cfg.Database.Driver)
				}
				if cfg.Database.Host != "localhost" {
					t.Errorf("expected database host localhost, got %s", cfg.Database.Host)
				}
				if cfg.Database.Port != 3306 {
					t.Errorf("expected database port 3306, got %d", cfg.Database.Port)
				}
				if cfg.Database.Name != "warehouse" {
					t.Errorf("expected database name warehouse, got %s", cfg.Database.Name)
				}
				if cfg.Database.User != "root" {
					t.Errorf("expected database user root, got %s", cfg.Database.User)
				}
				if cfg.Database.Password != "" {
					t.Errorf("expected database password empty, got %s", cfg.Database.Password)
				}
				if cfg.Database.MaxOpenConns != 100 {
					t.Errorf("expected max open conns 100, got %d", cfg.Database.MaxOpenConns)
				}
				if cfg.Database.MaxIdleConns != 10 {
					t.Errorf("expected max idle conns 10, got %d", cfg.Database.MaxIdleConns)
				}
				if cfg.JWT.Secret != "your-secret-key-change-in-production" {
					t.Errorf("expected jwt secret, got %s", cfg.JWT.Secret)
				}
				if cfg.JWT.Expire != "24h" {
					t.Errorf("expected jwt expire 24h, got %s", cfg.JWT.Expire)
				}
				if cfg.Log.Level != "debug" {
					t.Errorf("expected log level debug, got %s", cfg.Log.Level)
				}
				if cfg.Log.Output != "stdout" {
					t.Errorf("expected log output stdout, got %s", cfg.Log.Output)
				}
				if cfg.Log.File != "" {
					t.Errorf("expected log file empty, got %s", cfg.Log.File)
				}
			},
		},
		{
			name: "load non-existent file",
			setup: func() string {
				return "non-existent-config.yaml"
			},
			expectError: true,
		},
		{
			name: "load invalid yaml",
			setup: func() string {
				tmpFile := filepath.Join(os.TempDir(), "invalid-config.yaml")
				content := `
server:
  port: not-a-number
`
				_ = os.WriteFile(tmpFile, []byte(content), 0644)
				return tmpFile
			},
			expectError: true,
		},
		{
			name: "load partial config",
			setup: func() string {
				tmpFile := filepath.Join(os.TempDir(), "partial-config.yaml")
				content := `
server:
  port: 9090
`
				_ = os.WriteFile(tmpFile, []byte(content), 0644)
				return tmpFile
			},
			expectError: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Server.Port != 9090 {
					t.Errorf("expected server port 9090, got %d", cfg.Server.Port)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			cfg, err := Load(path)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if cfg == nil {
				t.Error("expected config, got nil")
				return
			}

			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}

func TestConfigDefaults(t *testing.T) {
	cfg, err := Load(filepath.Join("..", "..", "config", "config.yaml"))
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Server.Port == 0 {
		t.Error("server port should not be zero")
	}
	if cfg.Database.Port == 0 {
		t.Error("database port should not be zero")
	}
}

func TestConfigStructure(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: 8080,
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Driver:       "mysql",
			Host:         "localhost",
			Port:         3306,
			Name:         "warehouse",
			User:         "root",
			Password:     "",
			MaxOpenConns: 100,
			MaxIdleConns: 10,
		},
		JWT: JWTConfig{
			Secret: "test-secret",
			Expire: "24h",
		},
		Log: LogConfig{
			Level:  "debug",
			Output: "stdout",
			File:   "",
		},
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("expected server port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Database.Driver != "mysql" {
		t.Errorf("expected database driver mysql, got %s", cfg.Database.Driver)
	}
	if cfg.JWT.Secret != "test-secret" {
		t.Errorf("expected jwt secret test-secret, got %s", cfg.JWT.Secret)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("expected log level debug, got %s", cfg.Log.Level)
	}
}
