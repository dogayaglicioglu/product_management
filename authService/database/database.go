package database

import (
	"auth-service/models"
	"fmt"
	"log"
	"os"
	"time"

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
		"host=dbauth user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	fmt.Printf("DSN: %s\n", dsn)

	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			log.Println("Connected to the database successfully.")
			break
		}

		log.Printf("Failed to connect to database. Retrying in 5 seconds... (%d/10)\n", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to database after several attempts. \n", err)
		os.Exit(1)
	}
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Database migrations..")
	db.AutoMigrate(&models.AuthUser{})

	DB = DbInstance{
		DB: db,
	}
	return DB
}
