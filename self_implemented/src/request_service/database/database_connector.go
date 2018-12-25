package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var host string
var port int
var user string
var password string
var dbName string

var db sql.DB

type TaskStatus struct {
	ID string
	Status string
}

func Init(
	hostArg string,
	portArg int,
	userArg string,
	passwordArg string,
	dbNameArg string) error {

	host = hostArg
	port = portArg
	user = userArg
	password = passwordArg
	dbName = dbNameArg

	return exec(`
		CREATE TABLE IF NOT EXISTS Tasks (
			id VARCHAR(255) PRIMARY KEY,
			status VARCHAR(255)
		)
	`)
}

func Persist(taskId string) error {
	return exec(`
		INSERT INTO Tasks (id, status)
		VALUES ('` + taskId + `', 'new')
	`)
}

func UpdateStatus(taskId string, newStatus string) error {
	return exec(` 
		UPDATE Tasks
			SET status = '` + newStatus + `'
			WHERE id = '` + taskId + `'
	`)
}

func FetchStatus(taskId string) (string, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return "", err
	}
	defer db.Close()

	row := db.QueryRow(`SELECT status FROM Tasks where id='` + taskId + `'`)

	var status string
	err = row.Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func FetchAll() (*[]TaskStatus, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, e := db.Query("SELECT id, status FROM Tasks")
	if e != nil {
		panic(e)
	}
	defer rows.Close()

	tasks := make([]TaskStatus, 0)

	for rows.Next() {
		var t TaskStatus
		err := rows.Scan(&t.ID, &t.Status)
		if err != nil {
			panic(err)
		}

		tasks = append(tasks, t)
	}

	return &tasks, nil
}



func exec(query string) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(query)
	return err
}


func Shutdown() {
	db.Close()
}