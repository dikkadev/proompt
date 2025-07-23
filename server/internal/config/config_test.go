package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid local database config",
			config: Config{
				Database: Database{
					Local: &LocalDatabase{
						Path:       "/tmp/test.db",
						Migrations: "/tmp/migrations",
					},
				},
				Storage: Storage{
					ReposDir: "/tmp/repos",
				},
				Server: Server{
					Host: "localhost",
					Port: 8080,
				},
			},
			wantErr: false,
		},
		{
			name: "valid turso database config",
			config: Config{
				Database: Database{
					Turso: &TursoDatabase{
						URL:   "https://test.turso.io",
						Token: "test-token",
					},
				},
				Storage: Storage{
					ReposDir: "/tmp/repos",
				},
				Server: Server{
					Host: "localhost",
					Port: 8080,
				},
			},
			wantErr: false,
		},
		{
			name: "both database types configured - should fail",
			config: Config{
				Database: Database{
					Local: &LocalDatabase{
						Path:       "/tmp/test.db",
						Migrations: "/tmp/migrations",
					},
					Turso: &TursoDatabase{
						URL:   "https://test.turso.io",
						Token: "test-token",
					},
				},
				Storage: Storage{
					ReposDir: "/tmp/repos",
				},
				Server: Server{
					Host: "localhost",
					Port: 8080,
				},
			},
			wantErr: true,
			errMsg:  "must configure exactly one database type",
		},
		{
			name: "no database configured - should fail",
			config: Config{
				Database: Database{},
				Storage: Storage{
					ReposDir: "/tmp/repos",
				},
				Server: Server{
					Host: "localhost",
					Port: 8080,
				},
			},
			wantErr: true,
			errMsg:  "must configure exactly one database type",
		},
		{
			name: "missing required fields - should fail",
			config: Config{
				Database: Database{
					Local: &LocalDatabase{
						Path:       "", // missing
						Migrations: "/tmp/migrations",
					},
				},
				Storage: Storage{
					ReposDir: "", // missing
				},
				Server: Server{
					Host: "", // missing
					Port: 8080,
				},
			},
			wantErr: true,
			errMsg:  "cannot be empty",
		},
		{
			name: "invalid port - should fail",
			config: Config{
				Database: Database{
					Local: &LocalDatabase{
						Path:       "/tmp/test.db",
						Migrations: "/tmp/migrations",
					},
				},
				Storage: Storage{
					ReposDir: "/tmp/repos",
				},
				Server: Server{
					Host: "localhost",
					Port: 70000, // invalid port
				},
			},
			wantErr: true,
			errMsg:  "must be at most",
		},
		{
			name: "invalid turso URL - should fail",
			config: Config{
				Database: Database{
					Turso: &TursoDatabase{
						URL:   "not-a-url", // invalid URL
						Token: "test-token",
					},
				},
				Storage: Storage{
					ReposDir: "/tmp/repos",
				},
				Server: Server{
					Host: "localhost",
					Port: 8080,
				},
			},
			wantErr: true,
			errMsg:  "must be a valid URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errMsg != "" && !containsString(err.Error(), tt.errMsg) {
					t.Errorf("expected error to contain %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}

func TestEnvironmentConfigParsing(t *testing.T) {
	tests := []struct {
		name         string
		xmlContent   string
		environment  string
		wantHost     string
		wantPort     int
		wantDBType   string
		wantLogLevel string
		wantError    bool
	}{
		{
			name: "dev environment selection",
			xmlContent: `<?xml version="1.0" encoding="UTF-8"?>
<proompt>
    <database environment="dev">
        <local path="./dev.db" migrations="./migrations" />
    </database>
    <database environment="prod">
        <local path="/prod.db" migrations="/migrations" />
    </database>
    <storage environment="dev" repos_dir="./dev-repos" />
    <storage environment="prod" repos_dir="/prod-repos" />
    <server environment="dev" host="localhost" port="8080" />
    <server environment="prod" host="0.0.0.0" port="80" />
    <logging environment="dev" level="debug" />
    <logging environment="prod" level="info" />
</proompt>`,
			environment:  "dev",
			wantHost:     "localhost",
			wantPort:     8080,
			wantDBType:   "local",
			wantLogLevel: "debug",
			wantError:    false,
		},
		{
			name: "prod environment selection",
			xmlContent: `<?xml version="1.0" encoding="UTF-8"?>
<proompt>
    <database environment="dev">
        <local path="./dev.db" migrations="./migrations" />
    </database>
    <database environment="prod">
        <local path="/prod.db" migrations="/migrations" />
    </database>
    <storage environment="dev" repos_dir="./dev-repos" />
    <storage environment="prod" repos_dir="/prod-repos" />
    <server environment="dev" host="localhost" port="8080" />
    <server environment="prod" host="0.0.0.0" port="80" />
    <logging environment="dev" level="debug" />
    <logging environment="prod" level="info" />
</proompt>`,
			environment:  "prod",
			wantHost:     "0.0.0.0",
			wantPort:     80,
			wantDBType:   "local",
			wantLogLevel: "info",
			wantError:    false,
		},
		{
			name: "fallback to default (no environment attribute)",
			xmlContent: `<?xml version="1.0" encoding="UTF-8"?>
<proompt>
    <database>
        <local path="./default.db" migrations="./migrations" />
    </database>
    <storage repos_dir="./default-repos" />
    <server host="127.0.0.1" port="9000" />
    <logging level="warn" />
</proompt>`,
			environment:  "dev",
			wantHost:     "127.0.0.1",
			wantPort:     9000,
			wantDBType:   "local",
			wantLogLevel: "warn",
			wantError:    false,
		},
		{
			name: "mixed environment and default configs",
			xmlContent: `<?xml version="1.0" encoding="UTF-8"?>
<proompt>
    <database environment="prod">
        <local path="/prod.db" migrations="/migrations" />
    </database>
    <database>
        <local path="./default.db" migrations="./migrations" />
    </database>
    <storage repos_dir="./default-repos" />
    <server environment="dev" host="localhost" port="8080" />
    <server host="127.0.0.1" port="9000" />
    <logging level="info" />
</proompt>`,
			environment:  "dev",
			wantHost:     "localhost",
			wantPort:     8080,
			wantDBType:   "local",
			wantLogLevel: "info",
			wantError:    false,
		},
		{
			name: "missing environment config should error",
			xmlContent: `<?xml version="1.0" encoding="UTF-8"?>
<proompt>
    <database environment="prod">
        <local path="/prod.db" migrations="/migrations" />
    </database>
    <storage environment="prod" repos_dir="/prod-repos" />
    <server environment="prod" host="0.0.0.0" port="80" />
</proompt>`,
			environment: "dev",
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "test.xml")

			err := os.WriteFile(configPath, []byte(tt.xmlContent), 0644)
			if err != nil {
				t.Fatalf("Failed to write test config: %v", err)
			}

			// Load config
			config, err := Load(configPath, tt.environment)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify server config
			if config.Server.Host != tt.wantHost {
				t.Errorf("Server.Host = %v, want %v", config.Server.Host, tt.wantHost)
			}
			if config.Server.Port != tt.wantPort {
				t.Errorf("Server.Port = %v, want %v", config.Server.Port, tt.wantPort)
			}

			// Verify database type
			if config.DatabaseType() != tt.wantDBType {
				t.Errorf("DatabaseType() = %v, want %v", config.DatabaseType(), tt.wantDBType)
			}

			// Verify logging level
			if config.Logging.Level != tt.wantLogLevel {
				t.Errorf("Logging.Level = %v, want %v", config.Logging.Level, tt.wantLogLevel)
			}
		})
	}
}

func TestLoggingDefaults(t *testing.T) {
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<proompt>
    <database>
        <local path="./test.db" migrations="./migrations" />
    </database>
    <storage repos_dir="./repos" />
    <server host="localhost" port="8080" />
</proompt>`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.xml")

	err := os.WriteFile(configPath, []byte(xmlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	config, err := Load(configPath, "dev")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check logging defaults
	if config.Logging.Level != "debug" {
		t.Errorf("Default logging level = %v, want debug", config.Logging.Level)
	}
	if !config.Logging.Outputs.Stdout.Enabled {
		t.Errorf("Default stdout should be enabled")
	}
	if !config.Logging.Outputs.Stdout.Colors {
		t.Errorf("Default stdout colors should be enabled")
	}
	if config.Logging.Outputs.File.Enabled {
		t.Errorf("Default file logging should be disabled")
	}
}
