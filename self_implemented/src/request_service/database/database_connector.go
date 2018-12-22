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
	dbNameArg string) {

	host = hostArg
	port = portArg
	user = userArg
	password = passwordArg
	dbName = dbNameArg

	exec(`
		CREATE TABLE IF NOT EXISTS Tasks (
			id VARCHAR(255) PRIMARY KEY,
			status VARCHAR(255)
		)
	`)
}

func Persist(taskId string) {
	exec(`
		INSERT INTO Tasks (id, status)
		VALUES ('` + taskId + `', 'new')
	`)
}

func UpdateStatus(taskId string, newStatus string) {
	exec(` 
		UPDATE Tasks
			SET status = '` + newStatus + `'
			WHERE id = '` + taskId + `'
	`)
}

func FetchStatus(taskId string) string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	row := db.QueryRow(`SELECT status FROM Tasks where id='` + taskId + `'`)

	var status string
	err = row.Scan(&status)
	if err != nil {
		panic(err)
	}

	return status
}

func FetchAll() *[]TaskStatus {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
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

	return &tasks
}



func exec(query string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec(query)
}


func Shutdown() {
	db.Close()
}