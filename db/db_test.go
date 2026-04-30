package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name         string
		dbType       string
		host         string
		database     string
		expectedPort int
	}{
		{
			name:         "postgres config",
			dbType:       Postgres,
			host:         "localhost",
			database:     "testdb",
			expectedPort: 5432,
		},
		{
			name:         "mysql config",
			dbType:       MySQL,
			host:         "localhost",
			database:     "testdb",
			expectedPort: 3306,
		},
		{
			name:         "mssql config",
			dbType:       MSSQL,
			host:         "localhost",
			database:     "testdb",
			expectedPort: 1433,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig(tt.dbType, tt.host, tt.database)

			assert.Equal(t, tt.dbType, config.Type)
			assert.Equal(t, tt.host, config.Host)
			assert.Equal(t, tt.database, config.Database)
			assert.Equal(t, tt.expectedPort, config.Port)
			assert.Equal(t, 25, config.MaxOpenConns)
			assert.Equal(t, 5, config.MaxIdleConns)
			assert.True(t, config.EnableTracing)
			assert.True(t, config.EnableMetrics)
		})
	}
}

func TestConfigFluentBuilder(t *testing.T) {
	config := NewConfig(Postgres, "localhost", "testdb").
		WithPort(5433).
		WithCredentials("user", "pass").
		WithSSL("require", "/cert.pem", "/key.pem", "/root.pem").
		WithPool(50, 10, 10*time.Minute, 20*time.Minute).
		WithTracing(false).
		WithMetrics(false)

	assert.Equal(t, 5433, config.Port)
	assert.Equal(t, "user", config.Username)
	assert.Equal(t, "pass", config.Password)
	assert.Equal(t, "require", config.SSLMode)
	assert.Equal(t, "/cert.pem", config.SSLCert)
	assert.Equal(t, "/key.pem", config.SSLKey)
	assert.Equal(t, "/root.pem", config.SSLRootCert)
	assert.Equal(t, 50, config.MaxOpenConns)
	assert.Equal(t, 10, config.MaxIdleConns)
	assert.Equal(t, 10*time.Minute, config.ConnMaxLifetime)
	assert.Equal(t, 20*time.Minute, config.ConnMaxIdleTime)
	assert.False(t, config.EnableTracing)
	assert.False(t, config.EnableMetrics)
}

func TestConfigOptions(t *testing.T) {
	config := NewConfig(Postgres, "localhost", "testdb")

	// Apply options
	WithPort(5433)(config)
	WithCredentials("user", "pass")(config)
	WithSSL("require", "/cert.pem", "/key.pem", "/root.pem")(config)
	WithPool(50, 10, 10*time.Minute, 20*time.Minute)(config)
	WithTracing(false)(config)
	WithMetricsOption(false)(config)

	assert.Equal(t, 5433, config.Port)
	assert.Equal(t, "user", config.Username)
	assert.Equal(t, "pass", config.Password)
	assert.Equal(t, "require", config.SSLMode)
	assert.Equal(t, "/cert.pem", config.SSLCert)
	assert.Equal(t, "/key.pem", config.SSLKey)
	assert.Equal(t, "/root.pem", config.SSLRootCert)
	assert.Equal(t, 50, config.MaxOpenConns)
	assert.Equal(t, 10, config.MaxIdleConns)
	assert.Equal(t, 10*time.Minute, config.ConnMaxLifetime)
	assert.Equal(t, 20*time.Minute, config.ConnMaxIdleTime)
	assert.False(t, config.EnableTracing)
	assert.False(t, config.EnableMetrics)
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid postgres config",
			config: &Config{
				Type:     Postgres,
				Host:     "localhost",
				Database: "testdb",
			},
			wantErr: false,
		},
		{
			name: "invalid database type",
			config: &Config{
				Type:     "oracle",
				Host:     "localhost",
				Database: "testdb",
			},
			wantErr: true,
		},
		{
			name: "missing host",
			config: &Config{
				Type:     Postgres,
				Host:     "",
				Database: "testdb",
			},
			wantErr: true,
		},
		{
			name: "missing database",
			config: &Config{
				Type:     Postgres,
				Host:     "localhost",
				Database: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildDSN(t *testing.T) {
	t.Run("postgres DSN", func(t *testing.T) {
		config := &Config{
			Type:     Postgres,
			Host:     "localhost",
			Port:     5432,
			Database: "testdb",
			Username: "user",
			Password: "pass",
			SSLMode:  "disable",
		}

		dsn := buildPostgresDSN(config)
		assert.Contains(t, dsn, "host=localhost")
		assert.Contains(t, dsn, "port=5432")
		assert.Contains(t, dsn, "dbname=testdb")
		assert.Contains(t, dsn, "user=user")
		assert.Contains(t, dsn, "password=pass")
		assert.Contains(t, dsn, "sslmode=disable")
	})

	t.Run("mysql DSN", func(t *testing.T) {
		config := &Config{
			Type:     MySQL,
			Host:     "localhost",
			Port:     3306,
			Database: "testdb",
			Username: "user",
			Password: "pass",
			SSLMode:  "disable",
		}

		dsn := buildMySQLDSN(config)
		assert.Contains(t, dsn, "user:pass@tcp(localhost:3306)/testdb")
	})

	t.Run("mssql DSN", func(t *testing.T) {
		config := &Config{
			Type:     MSSQL,
			Host:     "localhost",
			Port:     1433,
			Database: "testdb",
			Username: "user",
			Password: "pass",
		}

		dsn := buildMSSQLDSN(config)
		assert.Contains(t, dsn, "sqlserver://user:pass@localhost:1433")
		assert.Contains(t, dsn, "database=testdb")
	})
}

func TestGetDefaultPort(t *testing.T) {
	assert.Equal(t, 5432, getDefaultPort(Postgres))
	assert.Equal(t, 3306, getDefaultPort(MySQL))
	assert.Equal(t, 1433, getDefaultPort(MSSQL))
	assert.Equal(t, 0, getDefaultPort("unknown"))
}

func TestBuildDialector(t *testing.T) {
	t.Run("postgres dialector", func(t *testing.T) {
		config := &Config{
			Type:     Postgres,
			Host:     "localhost",
			Port:     5432,
			Database: "testdb",
			Username: "user",
			Password: "pass",
		}
		dialector, err := buildDialector(config)
		require.NoError(t, err)
		assert.NotNil(t, dialector)
	})

	t.Run("mysql dialector", func(t *testing.T) {
		config := &Config{
			Type:     MySQL,
			Host:     "localhost",
			Port:     3306,
			Database: "testdb",
			Username: "user",
			Password: "pass",
		}
		dialector, err := buildDialector(config)
		require.NoError(t, err)
		assert.NotNil(t, dialector)
	})

	t.Run("mssql dialector", func(t *testing.T) {
		config := &Config{
			Type:     MSSQL,
			Host:     "localhost",
			Port:     1433,
			Database: "testdb",
			Username: "user",
			Password: "pass",
		}
		dialector, err := buildDialector(config)
		require.NoError(t, err)
		assert.NotNil(t, dialector)
	})

	t.Run("invalid database type", func(t *testing.T) {
		config := &Config{
			Type:     "oracle",
			Host:     "localhost",
			Port:     1521,
			Database: "testdb",
		}
		_, err := buildDialector(config)
		require.Error(t, err)
	})
}

func TestNewInvalidConfig(t *testing.T) {
	ctx := context.Background()

	// Test with invalid database type
	config := &Config{
		Type:     "invalid",
		Host:     "localhost",
		Database: "testdb",
	}

	_, err := New(ctx, config)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidDatabaseType)
}
