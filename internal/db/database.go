package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

var (
	DB   *gorm.DB
	once sync.Once
)

func GetDbConfig() *DbConfig {
	return &DbConfig{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "r00t"),
		DBName:     getEnv("DB_NAME", "edairy"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

// ConnectToDatabase opens the shared pool once for the process lifetime.
// Individual queries do not call Close; database/sql returns connections to the pool
// after each GORM operation completes.
func ConnectToDatabase() {
	once.Do(func() {
		cfg := GetDbConfig()

		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)

		logLevel := logger.Silent
		if getEnv("APP_ENV", "production") != "production" {
			logLevel = logger.Info
		}

		gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold:             time.Second,
					LogLevel:                  logLevel,
					IgnoreRecordNotFoundError: true,
					Colorful:                  true,
				},
			),
		})
		if err != nil {
			log.Fatal("Failed to connect database:", err)
		}

		sqlDB, err := gormDB.DB()
		if err != nil {
			log.Fatal("Failed to get db instance:", err)
		}

		configurePool(sqlDB)

		if err := sqlDB.Ping(); err != nil {
			log.Fatal("Failed to ping database:", err)
		}

		DB = gormDB
		log.Println("Database connected successfully")
	})
}

func configurePool(sqlDB interface {
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	SetConnMaxLifetime(d time.Duration)
	SetConnMaxIdleTime(d time.Duration)
}) {
	// Production MySQL wait_timeout is 600s — stay below that.
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(8 * time.Minute)
	sqlDB.SetConnMaxIdleTime(4 * time.Minute)
}

// WithContext binds the shared handle to ctx so work is cancelled when the
// request ends and the underlying connection is returned to the pool promptly.
func WithContext(ctx context.Context) *gorm.DB {
	if DB == nil {
		return nil
	}
	return DB.WithContext(ctx)
}

// CloseDatabase drains and closes the connection pool. Call on process shutdown.
func CloseDatabase() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	DB = nil
	log.Println("Database connection pool closed")
	return nil
}
