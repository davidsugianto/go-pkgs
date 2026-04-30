package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/davidsugianto/go-pkgs/db"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// User represents a sample model
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:255;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	ctx := context.Background()

	// Example 1: Simple connection with fluent config builder
	config := db.NewConfig(db.Postgres, "localhost", "myapp").
		WithPort(5432).
		WithCredentials("postgres", "password").
		WithPool(25, 5, 5*time.Minute, 10*time.Minute)

	client, err := db.New(ctx, config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Close()

	// Verify connection
	if err := client.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("Successfully connected to database!")

	// Example 2: Using options pattern (alternative approach)
	client2, err := db.New(ctx,
		db.NewConfig(db.MySQL, "localhost", "myapp"),
		db.WithPort(3306),
		db.WithCredentials("root", "password"),
		db.WithTracing(true),
		db.WithMetricsOption(true),
	)
	if err != nil {
		log.Printf("MySQL connection failed: %v", err)
	} else {
		defer client2.Close()
		fmt.Println("MySQL connection established!")
	}

	// Example 3: Auto-migrate models
	if err := client.DB.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
	fmt.Println("Migration completed!")

	// Example 4: Transaction
	err = client.Transaction(ctx, func(tx *gorm.DB) error {
		user := User{
			Name:  "John Doe",
			Email: "john@example.com",
		}
		return tx.Create(&user).Error
	})
	if err != nil {
		log.Printf("Transaction failed: %v", err)
	}

	// Example 5: Query with context
	var users []User
	if err := client.WithContext(ctx).Find(&users).Error; err != nil {
		log.Printf("Query failed: %v", err)
	}
	fmt.Printf("Found %d users\n", len(users))

	// Example 6: Connection pool stats
	stats, err := client.Stats()
	if err != nil {
		log.Printf("Failed to get stats: %v", err)
	} else {
		fmt.Printf("Connection pool: OpenConnections=%d, InUse=%d, Idle=%d\n",
			stats.OpenConnections, stats.InUse, stats.Idle)
	}

	// Example 7: Run embedded migrations
	// Uncomment to run migrations from embedded files
	// if err := client.Migrate(ctx, migrationsFS, "migrations"); err != nil {
	// 	log.Printf("Migration failed: %v", err)
	// }

	// Example 8: Health check
	if client.IsHealthy(ctx) {
		fmt.Println("Database is healthy!")
	}
}
