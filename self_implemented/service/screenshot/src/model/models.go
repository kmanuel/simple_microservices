package model

type Task struct {
	ID  string `jsonapi:"primary,screenshot_task"`
	Url string `jsonapi:"attr,url"`
}
