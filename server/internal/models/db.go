package models

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	gormDB *gorm.DB
	once   sync.Once
)

func Db() *gorm.DB {
	once.Do(func() {
		var err error

		env := os.Getenv("GO_ENV")
		log.Println("GO_ENV:", env)

		if env == "development" {
			err = godotenv.Load()
			if err != nil {
				log.Fatalf("Error loading .env file")
			}
			log.Println("Loaded .env var file")
		}

		log.Println("DB DSN: ", os.Getenv("HACKERNOON_DEV_DSN"))
		log.Println("DB DSN: ", os.Getenv("HACKERNOON_PROD_DSN"))

		switch env {
		case "development":
			// gormDB, err = gorm.Open(postgres.Open(os.Getenv("HACKERNOON_DEV_DSN")), &gorm.Config{
			gormDB, err = gorm.Open(postgres.Open(os.Getenv("HACKERNOON_PROD_DSN")), &gorm.Config{
				SkipDefaultTransaction: true,
				PrepareStmt:            true,
				Logger:                 logger.Default.LogMode(logger.Info),
			})
			if err != nil {
				log.Fatal("Failed to connect to the database: ", err)
			}
		case "production":
			gormDB, err = gorm.Open(postgres.Open(os.Getenv("HACKERNOON_PROD_DSN")), &gorm.Config{
				SkipDefaultTransaction: true, PrepareStmt: true,
			})
			if err != nil {
				log.Fatal("Failed to connect to the database: ", err)
			}
		default:
			log.Fatal("Unrecognized GO_ENV:", env)
		}

		log.Println("Connected to postgres successfully")

		err = gormDB.AutoMigrate(&Article{}, &Author{})
		if err != nil {
			log.Fatal("Failed to make auto migration", err)
		}
		log.Println("Auto Migration successful")
	})

	return gormDB
}
