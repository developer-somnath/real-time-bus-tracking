package models

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"real-time-bus-tracking/internal/utils"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var DB *gorm.DB

// RandomPolicy implements a random load-balancing policy for replicas
type RandomPolicy struct{}

func (RandomPolicy) Resolve(connPools []gorm.ConnPool) gorm.ConnPool {
	// Randomly select a replica from the pool
	selected := connPools[rand.Intn(len(connPools))]
	// Log the selected replica (could print its connection string or a custom identifier)
	fmt.Printf("Selecting replica database: %v\n", selected)
	return selected
}

// It sets the global DB variable to the opened database connection.
func InitDB() error {
	config := utils.LoadDBConfig()
	var err error
	for attempt := 1; attempt <= config.MaxRetries; attempt++ {
		DB, err = gorm.Open(mysql.Open(config.MasterDSN), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger:                 logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "t_", // Global prefix for all tables
			},
		})
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
			defer cancel()
			sqlDB, err := DB.DB()
			if err == nil && sqlDB.PingContext(ctx) == nil {
				break
			}
			if sqlDB != nil {
				sqlDB.Close()
			}
		}
		// Check if error is transient
		if utils.IsTransientError(err) {
			log.Printf("MySQL connection attempt %d/%d failed: %v", attempt, config.MaxRetries, err)
			if attempt < config.MaxRetries {
				// Add jitter to backoff
				jitter := time.Duration(rand.Intn(100)) * time.Millisecond
				time.Sleep(config.RetryBackoff*time.Duration(1<<uint(attempt-1)) + jitter)
			}
		} else {
			return fmt.Errorf("non-transient MySQL connection error: %w", err)
		}
	}
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL after %d retries: %w", config.MaxRetries, err)
	}
	err = DB.Use(dbresolver.Register(dbresolver.Config{
		Sources: []gorm.Dialector{mysql.Open(config.MasterDSN)},
		Replicas: func() []gorm.Dialector {
			dialectors := make([]gorm.Dialector, len(config.SlaveDSNs))
			for i, dsn := range config.SlaveDSNs {
				dialectors[i] = mysql.Open(dsn)
			}
			return dialectors
		}(),
		Policy:            RandomPolicy{},
		TraceResolverMode: true,
	}))
	if err != nil {
		return fmt.Errorf("failed to configure DB Resolver: %w", err)
	}
	// Set the global DB variable to the opened database connection
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sqlDB: %w", err)
	}
	// Connection pool settings
	sqlDB.SetMaxOpenConns(100)                 // Maximum number of open connections
	sqlDB.SetMaxIdleConns(10)                  // Maximum number of idle connections
	sqlDB.SetConnMaxLifetime(time.Hour)        // Maximum time a connection can be reused
	sqlDB.SetConnMaxIdleTime(30 * time.Minute) // Maximum time a connection can remain idle

	// Auto-migrate the schema
	// if os.Getenv("RUN_MIGRATE") == "true" {
	// 	err = DB.AutoMigrate(&schemas.User{}, &schemas.Order{})
	// 	if err != nil {
	// 		log.Printf("Failed to auto-migrate tables: %v", err)
	// 	} else {
	// 		fmt.Println("Tables auto-migrated successfully")
	// 	}
	// }
	return nil
}
