package database

import (
	"fmt"
	"log"
	"os"
	"product_management/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	DB *gorm.DB
}

var DB DbInstance

func ConnectDb() DbInstance {
	dsn := fmt.Sprintf(
		"host=dbweb user=%s password=%s dbname=%s port=5433 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	fmt.Printf("DSN: %s\n", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(1)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Database migrations..")
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Product{})

	DB = DbInstance{
		DB: db,
	}
	return DB
}
