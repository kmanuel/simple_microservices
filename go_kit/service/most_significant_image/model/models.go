package model

type MostSignificantImageTask struct {
	ID  string `jsonapi:"primary,most_significant_image_task"`
	Url string `jsonapi:"attr,url"`
}
