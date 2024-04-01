package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/JubaerHossain/golang-htmx-starter/pkg/core/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB *gorm.DB

var dbClient *gorm.DB

func GetDb() *gorm.DB {
	return dbClient
}

func ConnectDB() (DB, error) {
	// Create database connection
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=Asia/Dhaka",
		config.GlobalConfig.DBType,
		config.GlobalConfig.DBUser,
		config.GlobalConfig.DBPassword,
		config.GlobalConfig.DBHost,
		strconv.Itoa(config.GlobalConfig.DBPort),
		config.GlobalConfig.DBName,
	)

	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Disable prepared statements
	}), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		QueryFields: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection pool: %w", err)
	}
	sqlDB.SetMaxIdleConns(500)                  // Maximum number of idle connections in the pool
	sqlDB.SetMaxOpenConns(2000)                 // Maximum number of open connections to the database
	sqlDB.SetConnMaxLifetime(900 * time.Minute) // Maximum amount of time a connection may be reused
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)  // Maximum amount of time a connection may remain idle
	sqlDB.SetConnMaxIdleTime(10 * time.Second)  // Time period after which idle connections are closed

	// Test the database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i := 0; i < 3; i++ { // Retry logic
		err := sqlDB.PingContext(ctx)
		if err == nil {
			break
		}
		log.Printf("failed to ping database: %v (attempt %d)", err, i+1)
		time.Sleep(2 * time.Second) // Wait before retrying
	}
	if err != nil {
		return nil, fmt.Errorf("failed to ping database after multiple attempts: %w", err)
	}

	// defer func() {
	// 	dbInstance, _ := db.DB()
	// 	_ = dbInstance.Close()
	// }()

	log.Println("connected to database")

	return db, nil
}
func MigrateDB(db *gorm.DB) error {
	// Add your database migration logic here
	// For example:
	if err := db.AutoMigrate(
	); err != nil {
		return fmt.Errorf("failed to perform database migrations: %w", err)
	}

	log.Println("database migration completed")
	return nil
}
