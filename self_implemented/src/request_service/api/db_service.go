package api

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
)

func InitDb() {
	db, err := OpenDb()
	defer db.Close()
	if err != nil {
		panic(err)
		return
	}
	db.AutoMigrate(&TaskStatus{})
}

func OpenDb() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbName := os.Getenv("POSTGRES_DB")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := gorm.Open("postgres", psqlInfo)
	return db, err
}

func GetCountOfNotCompletedTasksOfType(taskType string) float64 {
	db, e := OpenDb()
	defer db.Close()
	if e != nil {
		return -1
	}

	var count float64
	db.Model(&TaskStatus{}).Where("task_type = ? AND status <> 'completed'", taskType).Count(&count)

	return count
}

