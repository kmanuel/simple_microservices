package model

type Task struct {
	ID      string `jsonapi:"primary,optimization_task"`
	ImageId string `jsonapi:"attr,image_id"`
	Width   int    `jsonapi:"attr,width"`
	Height  int    `jsonapi:"attr,height"`
}
