package model

type CropTask struct {
	ID      string `jsonapi:"primary,crop_task"`
	ImageId string `jsonapi:"attr,image_id"`
	Width   int    `jsonapi:"attr,width"`
	Height  int    `jsonapi:"attr,height"`
}
