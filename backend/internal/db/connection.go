package db

import (
    "log"
    "os"

    "github.com/zalg2261/bioskop/backend/internal/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

func Connect() {
    // Get database configuration from environment variables
    host := getEnv("DB_HOST", "localhost")
    user := getEnv("DB_USER", "postgres")
    password := getEnv("DB_PASSWORD", "0805")
    dbname := getEnv("DB_NAME", "db_bioskop")
    port := getEnv("DB_PORT", "2261") // Default port sesuai konfigurasi user

    dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"

    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        DisableForeignKeyConstraintWhenMigrating: true,
    })
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }
    
    // Try to drop existing constraints if they exist (ignore errors)
    // This prevents GORM from trying to drop non-existent constraints
    database.Exec("DO $$ BEGIN " +
        "ALTER TABLE IF EXISTS users DROP CONSTRAINT IF EXISTS uni_users_email; " +
        "EXCEPTION WHEN undefined_table THEN NULL; " +
        "END $$;")
    database.Exec("DO $$ BEGIN " +
        "ALTER TABLE IF EXISTS users DROP CONSTRAINT IF EXISTS users_email_key; " +
        "EXCEPTION WHEN undefined_table THEN NULL; " +
        "END $$;")
    
    err = database.AutoMigrate(
        &models.User{},
        &models.Movie{},
        &models.Showtime{},
        &models.Booking{},
        &models.Transaction{},
        &models.Wallet{},
        &models.City{},
        &models.Branch{},
        &models.Refund{},
        &models.SeatLock{},
    )

    if err != nil {
        log.Fatal("Failed to migrate:", err)
    }

    DB = database
    log.Println("DB connected & migrated!")
}
