package database

import (
	"backend/config"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) *gorm.DB {
	logMode := logger.Silent
	if os.Getenv("APP_ENV") == "development" {
		logMode = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logMode,
			},
		),
	})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan koneksi DB instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database tidak dapat diakses: %v", err)
	}

	log.Println("Koneksi database berhasil!")
	return db
}
