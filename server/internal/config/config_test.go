package config

import (
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
