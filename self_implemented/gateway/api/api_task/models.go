package api_task

type FaktoryInfo struct {
	ID             int                    `jsonapi:"primary,resources"`
	TotalProcessed float64                `jsonapi:"attr,totalProcessed"`
	TotalQueues    float64                `jsonapi:"attr,totalQueues"`
	TotalEnqueued  float64                `jsonapi:"attr,totalEnqueued"`
	TotalFailures  float64                `jsonapi:"attr,totalFailures"`
	Queues         map[string]float64     `jsonapi:"attr,queues"`
	Tasks          map[string]interface{} `jsonapi:"attr,tasks"`
}
