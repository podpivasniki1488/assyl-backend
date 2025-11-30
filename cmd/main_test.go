package main

import (
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: Config{
				RedisDSN:      "redis://localhost:6379",
				RedisUsername: "redis_user",
				RedisPassword: "redis_pass",
				DBDSN:         "postgres://user:pass@localhost/db",
				JwtSecretKey:  "secret-key-123",
				Debug:         false,
				GmailUsername: "test@gmail.com",
				GmailPassword: "gmail_pass",
				HttpPort:      "8080",
			},
			wantErr: false,
		},
		{
			name: "missing RedisDSN",
			config: Config{
				RedisDSN:      "",
				RedisUsername: "redis_user",
				RedisPassword: "redis_pass",
				DBDSN:         "postgres://user:pass@localhost/db",
				JwtSecretKey:  "secret-key-123",
				GmailUsername: "test@gmail.com",
				GmailPassword: "gmail_pass",
				HttpPort:      "8080",
			},
			wantErr: true,
			errMsg:  "RedisDSN",
		},
		{
			name: "missing RedisUsername",
			config: Config{
				RedisDSN:      "redis://localhost:6379",
				RedisUsername: "",
				RedisPassword: "redis_pass",
				DBDSN:         "postgres://user:pass@localhost/db",
				JwtSecretKey:  "secret-key-123",
				GmailUsername: "test@gmail.com",
				GmailPassword: "gmail_pass",
				HttpPort:      "8080",
			},
			wantErr: true,
			errMsg:  "RedisUsername",
		},
		{
			name: "missing RedisPassword",
			config: Config{
				RedisDSN:      "redis://localhost:6379",
				RedisUsername: "redis_user",
				RedisPassword: "",
				DBDSN:         "postgres://user:pass@localhost/db",
				JwtSecretKey:  "secret-key-123",
				GmailUsername: "test@gmail.com",
				GmailPassword: "gmail_pass",
				HttpPort:      "8080",
			},
			wantErr: true,
			errMsg:  "RedisPassword",
		},
		{
			name: "missing DBDSN",
			config: Config{
				RedisDSN:      "redis://localhost:6379",
				RedisUsername: "redis_user",
				RedisPassword: "redis_pass",
				DBDSN:         "",
				JwtSecretKey:  "secret-key-123",
				GmailUsername: "test@gmail.com",
				GmailPassword: "gmail_pass",
				HttpPort:      "8080",
			},
			wantErr: true,
			errMsg:  "DBDSN",
		},
		{
			name: "missing JwtSecretKey",
			config: Config{
				RedisDSN:      "redis://localhost:6379",
				RedisUsername: "redis_user",
				RedisPassword: "redis_pass",
				DBDSN:         "postgres://user:pass@localhost/db",
				JwtSecretKey:  "",
				GmailUsername: "test@gmail.com",
				GmailPassword: "gmail_pass",
				HttpPort:      "8080",
			},
			wantErr: true,
			errMsg:  "JwtSecretKey",
		},
		{
			name: "invalid email format for GmailUsername",
			config: Config{
				RedisDSN:      "redis://localhost:6379",
				RedisUsername: "redis_user",
				RedisPassword: "redis_pass",
				DBDSN:         "postgres://user:pass@localhost/db",
				JwtSecretKey:  "secret-key-123",
				GmailUsername: "not-an-email",
				GmailPassword: "gmail_pass",
				HttpPort:      "8080",
			},
			wantErr: true,
			errMsg:  "email",
		},
		{
			name: "missing GmailUsername",
			config: Config{
				RedisDSN:      "redis://localhost:6379",
				RedisUsername: "redis_user",
				RedisPassword: "redis_pass",
				DBDSN:         "postgres://user:pass@localhost/db",
				JwtSecretKey:  "secret-key-123",
				GmailUsername: "",
				GmailPassword: "gmail_pass",
				HttpPort:      "8080",
			},
			wantErr: true,
			errMsg:  "GmailUsername",
		},
		{
			name: "missing GmailPassword",
			config: Config{
				RedisDSN:      "redis://localhost:6379",
				RedisUsername: "redis_user",
				RedisPassword: "redis_pass",
				DBDSN:         "postgres://user:pass@localhost/db",
				JwtSecretKey:  "secret-key-123",
				GmailUsername: "test@gmail.com",
				GmailPassword: "",
				HttpPort:      "8080",
			},
			wantErr: true,
			errMsg:  "GmailPassword",
		},
		{
			name: "missing HttpPort",
			config: Config{
				RedisDSN:      "redis://localhost:6379",
				RedisUsername: "redis_user",
				RedisPassword: "redis_pass",
				DBDSN:         "postgres://user:pass@localhost/db",
				JwtSecretKey:  "secret-key-123",
				GmailUsername: "test@gmail.com",
				GmailPassword: "gmail_pass",
				HttpPort:      "",
			},
			wantErr: true,
			errMsg:  "HttpPort",
		},
		{
			name: "debug flag true",
			config: Config{
				RedisDSN:      "redis://localhost:6379",
				RedisUsername: "redis_user",
				RedisPassword: "redis_pass",
				DBDSN:         "postgres://user:pass@localhost/db",
				JwtSecretKey:  "secret-key-123",
				Debug:         true,
				GmailUsername: "test@gmail.com",
				GmailPassword: "gmail_pass",
				HttpPort:      "8080",
			},
			wantErr: false,
		},
		{
			name: "debug flag false",
			config: Config{
				RedisDSN:      "redis://localhost:6379",
				RedisUsername: "redis_user",
				RedisPassword: "redis_pass",
				DBDSN:         "postgres://user:pass@localhost/db",
				JwtSecretKey:  "secret-key-123",
				Debug:         false,
				GmailUsername: "test@gmail.com",
				GmailPassword: "gmail_pass",
				HttpPort:      "8080",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.New().Struct(&tt.config)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Config validation error = %v, wantErr %v", err, tt.wantErr)
			}
			
			if tt.wantErr && err != nil {
				errMsg := err.Error()
				if tt.errMsg != "" && len(errMsg) == 0 {
					t.Errorf("Expected error message to contain %q, got empty", tt.errMsg)
				}
			}
		})
	}
}

func TestMustReadConfig_Integration(t *testing.T) {
	// Save original env vars
	originalEnvVars := map[string]string{
		"REDIS_DSN":      os.Getenv("REDIS_DSN"),
		"REDIS_USERNAME": os.Getenv("REDIS_USERNAME"),
		"REDIS_PASSWORD": os.Getenv("REDIS_PASSWORD"),
		"JWT_SECRET":     os.Getenv("JWT_SECRET"),
		"DB_DSN":         os.Getenv("DB_DSN"),
		"DEBUG":          os.Getenv("DEBUG"),
		"GMAIL_USERNAME": os.Getenv("GMAIL_USERNAME"),
		"GMAIL_PASSWORD": os.Getenv("GMAIL_PASSWORD"),
		"PORT":           os.Getenv("PORT"),
	}
	
	// Cleanup after test
	defer func() {
		for key, val := range originalEnvVars {
			if val == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, val)
			}
		}
	}()
	
	t.Run("valid environment variables", func(t *testing.T) {
		os.Setenv("REDIS_DSN", "redis://localhost:6379")
		os.Setenv("REDIS_USERNAME", "redis_user")
		os.Setenv("REDIS_PASSWORD", "redis_pass")
		os.Setenv("JWT_SECRET", "jwt_secret")
		os.Setenv("DB_DSN", "postgres://user:pass@localhost/db")
		os.Setenv("DEBUG", "true")
		os.Setenv("GMAIL_USERNAME", "test@gmail.com")
		os.Setenv("GMAIL_PASSWORD", "gmail_pass")
		os.Setenv("PORT", "8080")
		
		// This should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("mustReadConfig() panicked with valid config: %v", r)
			}
		}()
		
		cfg := mustReadConfig()
		
		if cfg.RedisDSN != "redis://localhost:6379" {
			t.Errorf("RedisDSN = %v, want redis://localhost:6379", cfg.RedisDSN)
		}
		
		if cfg.Debug != true {
			t.Error("Debug should be true")
		}
		
		if cfg.HttpPort != "8080" {
			t.Errorf("HttpPort = %v, want 8080", cfg.HttpPort)
		}
	})
	
	t.Run("DEBUG=false", func(t *testing.T) {
		os.Setenv("REDIS_DSN", "redis://localhost:6379")
		os.Setenv("REDIS_USERNAME", "redis_user")
		os.Setenv("REDIS_PASSWORD", "redis_pass")
		os.Setenv("JWT_SECRET", "jwt_secret")
		os.Setenv("DB_DSN", "postgres://user:pass@localhost/db")
		os.Setenv("DEBUG", "false")
		os.Setenv("GMAIL_USERNAME", "test@gmail.com")
		os.Setenv("GMAIL_PASSWORD", "gmail_pass")
		os.Setenv("PORT", "8080")
		
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("mustReadConfig() panicked: %v", r)
			}
		}()
		
		cfg := mustReadConfig()
		
		if cfg.Debug != false {
			t.Error("Debug should be false")
		}
	})
	
	t.Run("missing required field should panic", func(t *testing.T) {
		os.Setenv("REDIS_DSN", "redis://localhost:6379")
		os.Setenv("REDIS_USERNAME", "redis_user")
		os.Setenv("REDIS_PASSWORD", "redis_pass")
		os.Setenv("JWT_SECRET", "jwt_secret")
		os.Setenv("DB_DSN", "postgres://user:pass@localhost/db")
		os.Setenv("DEBUG", "true")
		os.Setenv("GMAIL_USERNAME", "test@gmail.com")
		os.Setenv("GMAIL_PASSWORD", "gmail_pass")
		os.Unsetenv("PORT") // Missing required field
		
		defer func() {
			if r := recover(); r == nil {
				t.Error("mustReadConfig() should panic with missing PORT")
			}
		}()
		
		mustReadConfig()
	})
	
	t.Run("invalid email should panic", func(t *testing.T) {
		os.Setenv("REDIS_DSN", "redis://localhost:6379")
		os.Setenv("REDIS_USERNAME", "redis_user")
		os.Setenv("REDIS_PASSWORD", "redis_pass")
		os.Setenv("JWT_SECRET", "jwt_secret")
		os.Setenv("DB_DSN", "postgres://user:pass@localhost/db")
		os.Setenv("DEBUG", "true")
		os.Setenv("GMAIL_USERNAME", "not-an-email") // Invalid email
		os.Setenv("GMAIL_PASSWORD", "gmail_pass")
		os.Setenv("PORT", "8080")
		
		defer func() {
			if r := recover(); r == nil {
				t.Error("mustReadConfig() should panic with invalid email")
			}
		}()
		
		mustReadConfig()
	})
}

func TestConfig_StructTags(t *testing.T) {
	// Test that struct tags are properly defined
	v := validator.New()
	
	validCfg := Config{
		RedisDSN:      "redis://localhost:6379",
		RedisUsername: "user",
		RedisPassword: "pass",
		DBDSN:         "postgres://localhost/db",
		JwtSecretKey:  "secret",
		GmailUsername: "test@gmail.com",
		GmailPassword: "pass",
		HttpPort:      "8080",
	}
	
	if err := v.Struct(&validCfg); err != nil {
		t.Errorf("Valid config failed validation: %v", err)
	}
	
	// Test each validation tag
	testCases := []struct {
		name   string
		modify func(*Config)
	}{
		{
			name: "RedisDSN required",
			modify: func(c *Config) {
				c.RedisDSN = ""
			},
		},
		{
			name: "RedisUsername required",
			modify: func(c *Config) {
				c.RedisUsername = ""
			},
		},
		{
			name: "RedisPassword required",
			modify: func(c *Config) {
				c.RedisPassword = ""
			},
		},
		{
			name: "DBDSN required",
			modify: func(c *Config) {
				c.DBDSN = ""
			},
		},
		{
			name: "JwtSecretKey required",
			modify: func(c *Config) {
				c.JwtSecretKey = ""
			},
		},
		{
			name: "GmailUsername required and email",
			modify: func(c *Config) {
				c.GmailUsername = "invalid"
			},
		},
		{
			name: "GmailPassword required",
			modify: func(c *Config) {
				c.GmailPassword = ""
			},
		},
		{
			name: "HttpPort required",
			modify: func(c *Config) {
				c.HttpPort = ""
			},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := validCfg
			tc.modify(&cfg)
			
			err := v.Struct(&cfg)
			if err == nil {
				t.Error("Expected validation error, got nil")
			}
		})
	}
}