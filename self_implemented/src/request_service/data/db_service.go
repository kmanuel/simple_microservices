package data

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/kmanuel/simple_microservices/self_implemented/src/request_service/model"
	"github.com/prometheus/common/log"
	"os"
)

func InitDb() {
	db, err := OpenDb()
	defer closeDb(db)
	if err != nil {
		panic(err)
		return
	}
	db.AutoMigrate(&model.TaskStatus{})
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

func closeDb(db *gorm.DB) {
	err := db.Close()
	if err != nil {
		log.Error(err)
	}
}

func GetCountOfNotCompletedTasksOfType(taskType string) float64 {
	db, e := OpenDb()
	defer closeDb(db)
	if e != nil {
		return -1
	}

	var count float64
	db.Model(&model.TaskStatus{}).Where("task_type = ? AND status <> 'completed'", taskType).Count(&count)

	return count
}

func FetchTaskList() (*model.TaskStatusList, error) {
	db, err := OpenDb()
	defer closeDb(db)
	if err != nil {
		return nil, err
	}

	var tasks []*model.TaskStatus
	if err := db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	list := model.TaskStatusList{
		ID:    "1",
		Tasks: tasks,
	}
	return &list, nil
}

func CreateNewTask(newTask *model.TaskStatus) error {
	newTask.Status = "new"

	db, err := OpenDb()
	defer closeDb(db)
	if err != nil {
		return err
	}

	db.Create(&newTask)
	return nil
}

func UpdateTaskStatus(updateRequest *model.TaskStatus) error {
	db, err := OpenDb()
	defer closeDb(db)
	if err != nil {
		return err
	}

	var taskStatus model.TaskStatus
	if err := db.Where("task_id = ?", updateRequest.TaskID).First(&taskStatus).Error; err != nil {
		return err
	}

	taskStatus.Status = updateRequest.Status
	db.Save(&taskStatus)

	return nil
}
