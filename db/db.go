package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormotel "gorm.io/plugin/opentelemetry/tracing"
)

// Database type constants
const (
	Postgres = "postgres"
	MySQL    = "mysql"
	MSSQL    = "mssql"
)

var (
	ErrInvalidDatabaseType = errors.New("invalid database type")
	ErrConnectionFailed    = errors.New("connection failed")
	ErrMigrationFailed     = errors.New("migration failed")
)

// Config holds database connection configuration
type Config struct {
	// Database type (postgres, mysql, mssql)
	Type string `json:"type" yaml:"type" validate:"required,oneof=postgres mysql mssql"`

	// Connection settings
	Host     string `json:"host" yaml:"host" validate:"required"`
	Port     int    `json:"port" yaml:"port" validate:"required,min=1,max=65535"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database" validate:"required"`

	// SSL settings
	SSLMode     string `json:"ssl_mode" yaml:"ssl_mode"`           // for postgres: disable, require, verify-ca, verify-full
	SSLCert     string `json:"ssl_cert" yaml:"ssl_cert"`           // path to cert file
	SSLKey      string `json:"ssl_key" yaml:"ssl_key"`             // path to key file
	SSLRootCert string `json:"ssl_root_cert" yaml:"ssl_root_cert"` // path to root cert file

	// Connection pool settings
	MaxOpenConns    int           `json:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time" yaml:"conn_max_idle_time"`

	// OpenTelemetry settings
	EnableTracing bool `json:"enable_tracing" yaml:"enable_tracing"`
	EnableMetrics bool `json:"enable_metrics" yaml:"enable_metrics"`

	// GORM settings
	SkipDefaultTransaction bool `json:"skip_default_transaction" yaml:"skip_default_transaction"`
	PrepareStmt            bool `json:"prepare_stmt" yaml:"prepare_stmt"`
}

// Client wraps gorm.DB with additional functionality
type Client struct {
	*gorm.DB
	config     *Config
	meter      metric.Meter
	poolMetric metric.Int64ObservableGauge
}

// Option configures the database client
type Option func(*Config)

// NewConfig creates a new Config with defaults
func NewConfig(dbType, host, database string) *Config {
	return &Config{
		Type:            dbType,
		Host:            host,
		Port:            getDefaultPort(dbType),
		Database:        database,
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
		EnableTracing:   true,
		EnableMetrics:   true,
		PrepareStmt:     true,
	}
}

// WithPort sets the database port
func (c *Config) WithPort(port int) *Config {
	c.Port = port
	return c
}

// WithCredentials sets the username and password
func (c *Config) WithCredentials(username, password string) *Config {
	c.Username = username
	c.Password = password
	return c
}

// WithSSL sets SSL configuration
func (c *Config) WithSSL(mode, cert, key, rootCert string) *Config {
	c.SSLMode = mode
	c.SSLCert = cert
	c.SSLKey = key
	c.SSLRootCert = rootCert
	return c
}

// WithPool sets connection pool settings
func (c *Config) WithPool(maxOpen, maxIdle int, maxLifetime, maxIdleTime time.Duration) *Config {
	c.MaxOpenConns = maxOpen
	c.MaxIdleConns = maxIdle
	c.ConnMaxLifetime = maxLifetime
	c.ConnMaxIdleTime = maxIdleTime
	return c
}

// WithTracing enables or disables OpenTelemetry tracing
func (c *Config) WithTracing(enabled bool) *Config {
	c.EnableTracing = enabled
	return c
}

// WithMetrics enables or disables OpenTelemetry metrics
func (c *Config) WithMetrics(enabled bool) *Config {
	c.EnableMetrics = enabled
	return c
}

// WithPort sets the database port (as Option for functional options pattern)
func WithPort(port int) Option {
	return func(c *Config) {
		c.Port = port
	}
}

// WithCredentials sets the username and password (as Option for functional options pattern)
func WithCredentials(username, password string) Option {
	return func(c *Config) {
		c.Username = username
		c.Password = password
	}
}

// WithSSL sets SSL configuration (as Option for functional options pattern)
func WithSSL(mode, cert, key, rootCert string) Option {
	return func(c *Config) {
		c.SSLMode = mode
		c.SSLCert = cert
		c.SSLKey = key
		c.SSLRootCert = rootCert
	}
}

// WithPool sets connection pool settings (as Option for functional options pattern)
func WithPool(maxOpen, maxIdle int, maxLifetime, maxIdleTime time.Duration) Option {
	return func(c *Config) {
		c.MaxOpenConns = maxOpen
		c.MaxIdleConns = maxIdle
		c.ConnMaxLifetime = maxLifetime
		c.ConnMaxIdleTime = maxIdleTime
	}
}

// WithTracing enables or disables OpenTelemetry tracing (as Option for functional options pattern)
func WithTracing(enabled bool) Option {
	return func(c *Config) {
		c.EnableTracing = enabled
	}
}

// WithMetrics enables or disables OpenTelemetry metrics (as Option for functional options pattern)
func WithMetricsOption(enabled bool) Option {
	return func(c *Config) {
		c.EnableMetrics = enabled
	}
}

// New creates a new database client
func New(ctx context.Context, config *Config, opts ...Option) (*Client, error) {
	// Apply additional options
	for _, opt := range opts {
		opt(config)
	}

	// Validate config
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidDatabaseType, err)
	}

	// Build DSN and dialector
	dialector, err := buildDialector(config)
	if err != nil {
		return nil, err
	}

	// Build GORM config
	gormConfig := &gorm.Config{
		SkipDefaultTransaction: config.SkipDefaultTransaction,
		PrepareStmt:            config.PrepareStmt,
	}

	// Open connection
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	// Get underlying SQL DB for pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	client := &Client{
		DB:     db,
		config: config,
	}

	// Setup OpenTelemetry tracing
	if config.EnableTracing {
		if err := client.setupTracing(); err != nil {
			return nil, fmt.Errorf("failed to setup tracing: %w", err)
		}
	}

	// Setup OpenTelemetry metrics
	if config.EnableMetrics {
		if err := client.setupMetrics(ctx); err != nil {
			return nil, fmt.Errorf("failed to setup metrics: %w", err)
		}
	}

	return client, nil
}

// MustNew creates a new database client or panics on error
func MustNew(ctx context.Context, config *Config, opts ...Option) *Client {
	client, err := New(ctx, config, opts...)
	if err != nil {
		panic(err)
	}
	return client
}

// setupTracing configures OpenTelemetry tracing for GORM
func (c *Client) setupTracing() error {
	return c.DB.Use(gormotel.NewPlugin(
		gormotel.WithDBSystem(c.config.Type),
	))
}

// setupMetrics configures OpenTelemetry metrics for connection pool monitoring
func (c *Client) setupMetrics(ctx context.Context) error {
	c.meter = otel.GetMeterProvider().Meter("github.com/davidsugianto/go-pkgs/db")

	var err error
	c.poolMetric, err = c.meter.Int64ObservableGauge(
		"db.connection.pool.size",
		metric.WithDescription("Number of database connections in the pool"),
		metric.WithUnit("{connections}"),
	)
	if err != nil {
		return err
	}

	// Register callback for pool metrics
	_, err = c.meter.RegisterCallback(
		func(ctx context.Context, obs metric.Observer) error {
			if c.poolMetric == nil {
				return nil
			}

			sqlDB, err := c.DB.DB()
			if err != nil {
				return err
			}

			stats := sqlDB.Stats()

			attrs := []attribute.KeyValue{
				semconv.DBSystemKey.String(c.config.Type),
				semconv.DBNameKey.String(c.config.Database),
				attribute.String("db.host", c.config.Host),
			}

			obs.ObserveInt64(c.poolMetric, int64(stats.MaxOpenConnections), metric.WithAttributes(
				append(attrs, attribute.String("state", "max"))...,
			))
			obs.ObserveInt64(c.poolMetric, int64(stats.OpenConnections), metric.WithAttributes(
				append(attrs, attribute.String("state", "open"))...,
			))
			obs.ObserveInt64(c.poolMetric, int64(stats.InUse), metric.WithAttributes(
				append(attrs, attribute.String("state", "in_use"))...,
			))
			obs.ObserveInt64(c.poolMetric, int64(stats.Idle), metric.WithAttributes(
				append(attrs, attribute.String("state", "idle"))...,
			))
			obs.ObserveInt64(c.poolMetric, int64(stats.WaitCount), metric.WithAttributes(
				append(attrs, attribute.String("state", "wait_count"))...,
			))

			return nil
		},
		c.poolMetric,
	)
	return err
}

// Ping verifies the database connection
func (c *Client) Ping(ctx context.Context) error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Stats returns connection pool statistics
func (c *Client) Stats() (*sql.DBStats, error) {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return nil, err
	}
	stats := sqlDB.Stats()
	return &stats, nil
}

// Close closes the database connection
func (c *Client) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Migrate runs database migrations using embedded migration files
func (c *Client) Migrate(ctx context.Context, migrationsFS embed.FS, migrationsPath string) error {
	// Get underlying SQL DB
	sqlDB, err := c.DB.DB()
	if err != nil {
		return fmt.Errorf("%w: failed to get sql.DB: %v", ErrMigrationFailed, err)
	}

	// Create migrator using database-specific driver
	migrator, err := createMigrator(sqlDB, c.config, migrationsFS, migrationsPath)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMigrationFailed, err)
	}

	// Run migrations
	if err := migrator.Up(); err != nil {
		return fmt.Errorf("%w: %v", ErrMigrationFailed, err)
	}

	return nil
}

// Transaction executes a function within a database transaction
func (c *Client) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return c.DB.WithContext(ctx).Transaction(fn)
}

// WithContext returns a new client with the given context
func (c *Client) WithContext(ctx context.Context) *Client {
	return &Client{
		DB:     c.DB.WithContext(ctx),
		config: c.config,
		meter:  c.meter,
	}
}

// IsHealthy checks if the database connection is healthy
func (c *Client) IsHealthy(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return c.Ping(ctx) == nil
}

// getDefaultPort returns the default port for a database type
func getDefaultPort(dbType string) int {
	switch dbType {
	case Postgres:
		return 5432
	case MySQL:
		return 3306
	case MSSQL:
		return 1433
	default:
		return 0
	}
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	if config.Type != Postgres && config.Type != MySQL && config.Type != MSSQL {
		return fmt.Errorf("unsupported database type: %s (supported: %s, %s, %s)",
			config.Type, Postgres, MySQL, MSSQL)
	}
	if config.Host == "" {
		return errors.New("host is required")
	}
	if config.Database == "" {
		return errors.New("database name is required")
	}
	return nil
}

// buildDialector creates a GORM dialector based on database type
func buildDialector(config *Config) (gorm.Dialector, error) {
	switch config.Type {
	case Postgres:
		return postgres.Open(buildPostgresDSN(config)), nil
	case MySQL:
		return mysql.Open(buildMySQLDSN(config)), nil
	case MSSQL:
		return sqlserver.Open(buildMSSQLDSN(config)), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidDatabaseType, config.Type)
	}
}

// buildPostgresDSN builds a PostgreSQL connection string
func buildPostgresDSN(config *Config) string {
	dsn := fmt.Sprintf("host=%s port=%d dbname=%s",
		config.Host, config.Port, config.Database)

	if config.Username != "" {
		dsn += fmt.Sprintf(" user=%s", config.Username)
	}
	if config.Password != "" {
		dsn += fmt.Sprintf(" password=%s", config.Password)
	}
	if config.SSLMode != "" {
		dsn += fmt.Sprintf(" sslmode=%s", config.SSLMode)
	}
	if config.SSLCert != "" {
		dsn += fmt.Sprintf(" sslcert=%s", config.SSLCert)
	}
	if config.SSLKey != "" {
		dsn += fmt.Sprintf(" sslkey=%s", config.SSLKey)
	}
	if config.SSLRootCert != "" {
		dsn += fmt.Sprintf(" sslrootcert=%s", config.SSLRootCert)
	}

	return dsn
}

// buildMySQLDSN builds a MySQL connection string
func buildMySQLDSN(config *Config) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.Username, config.Password, config.Host, config.Port, config.Database)

	params := make([]string, 0)
	if config.SSLMode == "disable" {
		params = append(params, "tls=false")
	}
	if config.PrepareStmt {
		params = append(params, "interpolateParams=true")
	}

	if len(params) > 0 {
		dsn += "?"
		for i, p := range params {
			if i > 0 {
				dsn += "&"
			}
			dsn += p
		}
	}

	return dsn
}

// buildMSSQLDSN builds a MSSQL connection string
func buildMSSQLDSN(config *Config) string {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		config.Username, config.Password, config.Host, config.Port, config.Database)

	if config.SSLMode != "" && config.SSLMode != "disable" {
		dsn += fmt.Sprintf("&encrypt=true&TrustServerCertificate=%v",
			config.SSLMode == "require")
	}

	return dsn
}
