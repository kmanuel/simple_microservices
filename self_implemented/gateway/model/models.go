package model

type CropTask struct {
	ID      string `jsonapi:"primary,crop_task"`
	ImageId string `jsonapi:"attr,image_id"`
	Width   int    `jsonapi:"attr,width"`
	Height  int    `jsonapi:"attr,height"`
}

type MostSignificantImageTask struct {
	ID  string `jsonapi:"primary,most_significant_image_task"`
	Url string `jsonapi:"attr,url"`
}

type ScreenshotTask struct {
	ID  string `jsonapi:"primary,screenshot_task"`
	Url string `jsonapi:"attr,url"`
}

type OptimizationTask struct {
	ID      string `jsonapi:"primary,optimization_task"`
	ImageId string `jsonapi:"attr,image_id"`
	Width   int    `jsonapi:"attr,width"`
	Height  int    `jsonapi:"attr,height"`
}

type PortraitTask struct {
	ID      string `jsonapi:"primary,portrait_task"`
	ImageId string `jsonapi:"attr,image_id"`
	Width   int    `jsonapi:"attr,width"`
	Height  int    `jsonapi:"attr,height"`
}

type FaktoryInfo struct {
	ID             int                    `jsonapi:"primary,resources"`
	TotalProcessed float64                `jsonapi:"attr,totalProcessed"`
	TotalQueues    float64                `jsonapi:"attr,totalQueues"`
	TotalEnqueued  float64                `jsonapi:"attr,totalEnqueued"`
	TotalFailures  float64                `jsonapi:"attr,totalFailures"`
	Queues         map[string]float64     `jsonapi:"attr,queues"`
	Tasks          map[string]interface{} `jsonapi:"attr,tasks"`
}
