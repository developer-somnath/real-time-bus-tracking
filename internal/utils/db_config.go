package util

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	MaxRetries     int
	RetryBackoff   time.Duration
	ConnectTimeout time.Duration
	MasterDSN      string
	SlaveDSNs      []string
}

// loadDBConfig loads MySQL retry settings from environment variables
func LoadDBConfig() DBConfig {
	// dsn := "root:@tcp(localhost:3306)/gin_ai?charset=utf8mb4&parseTime=True&loc=Local"
	masterDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MASTER_DB_USER"), os.Getenv("MASTER_DB_PASS"), os.Getenv("MASTER_DB_HOST"), os.Getenv("MASTER_DB_PORT"), os.Getenv("MASTER_DB_NAME"))
	slave1Dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("SLAVE_1_DB_USER"), os.Getenv("SLAVE_1_DB_PASS"), os.Getenv("SLAVE_1_DB_HOST"), os.Getenv("SLAVE_1_DB_PORT"), os.Getenv("SLAVE_1_DB_NAME"))
	slave2Dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("SLAVE_2_DB_USER"), os.Getenv("SLAVE_2_DB_PASS"), os.Getenv("SLAVE_2_DB_HOST"), os.Getenv("SLAVE_2_DB_PORT"), os.Getenv("SLAVE_2_DB_NAME"))
	// InitDB initializes the database connection using GORM.
	return DBConfig{
		MaxRetries:     GetEnvInt("DB_MAX_RETRIES", 10),
		RetryBackoff:   GetEnvDuration("DB_RETRY_BACKOFF", 1000*time.Millisecond),
		ConnectTimeout: GetEnvDuration("DB_CONN_TIMEOUT", 5000*time.Millisecond),
		MasterDSN:      masterDsn,
		SlaveDSNs:      []string{slave1Dsn, slave2Dsn},
	}
}
func IsTransientError(err error) bool {
	if err == nil {
		return false
	}

	// Handle MySQL-specific errors
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok {
		// Common transient MySQL error codes
		if mysqlErr.Number == 1040 || mysqlErr.Number == 2003 || mysqlErr.Number == 2006 {
			return true
		}
	}

	// Handle network errors
	// We check if the error is a timeout error or connection failure
	if ne, ok := err.(*net.OpError); ok {
		// If the error is due to temporary network failure (e.g., connection refused or timeout), it might be transient
		if ne.Op == "dial" || ne.Op == "read" || ne.Op == "write" {
			// Assume it's transient (i.e., retryable)
			return true
		}
	}

	// Handle context timeout or cancellation errors
	if errors.Is(err, context.DeadlineExceeded) {
		// Timeout is generally a transient error, it can be retried
		return true
	}

	// Handle other network-related errors that may indicate a transient issue
	if errors.Is(err, context.Canceled) {
		// Canceled context might imply the operation was interrupted, it might be retryable
		return true
	}

	// Handle non-MySQL and non-network errors
	// For simplicity, assume any other error might be transient
	return true
}
