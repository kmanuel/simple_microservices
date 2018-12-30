package main

import (
	"bytes"
	"flag"
	"github.com/advancedlogic/GoOse"
	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/self_implemented/src/service/most_significant_image/update_status"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

const OutputImageLocation = "/tmp/"

type Task struct {
	ID  string `jsonapi:"primary,most_significant_image_task"`
	Url string `jsonapi:"attr,url"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	minioconnector.Init(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("BUCKET_NAME"))

	go startPrometheus()

	startFaktory()
}

var (
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Number of requests handled from faktory.",
		},
		[]string{"controller", "status"},
	)
)

func startPrometheus() {
	prometheus.MustRegister(requests)

	var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func startFaktory() {
	mgr := worker.NewManager()
	mgr.Use(func(perform worker.Handler) worker.Handler {
		return func(ctx worker.Context, job *faktory.Job) error {
			log.Printf("Starting work on job %s of type %s with custom %v\n", ctx.Jid(), ctx.JobType(), job.Custom)
			err := perform(ctx, job)
			log.Printf("Finished work on job %s with error %v\n", ctx.Jid(), err)
			return err
		}
	})
	mgr.Register("most_significant_image", convertTask)
	mgr.Queues = []string{"most_significant_image"}
	var quit bool
	mgr.On(worker.Shutdown, func() {
		quit = true
	})
	mgr.Run()
}

func convertTask(ctx worker.Context, args ...interface{}) error {
	log.Info("Working on job %s\n", ctx.Jid())

	task := new(Task)
	err := jsonapi.NewRuntime().UnmarshalPayload(bytes.NewBufferString(args[0].(string)), task)
	if err != nil {
		log.Error("failed to deserialize task", args)
		return err
	}

	update_status.NotifyAboutProcessingStart(task.ID)

	outputFile := OutputImageLocation + uuid.New().String() + ".jpg"

	err = ExtractMostSignificantImage(task.Url, outputFile)
	if err != nil {
		_ = ctx.Err()
		return nil
	}

	_, err = minioconnector.UploadFileWithName(outputFile, task.ID)
	if err != nil {
		_ = ctx.Err()
		return nil
	}

	update_status.NotifyAboutCompletion(task.ID)

	ctx.Done()
	return nil
}

func ExtractMostSignificantImage(inputUrl string, outputFile string) error {
	g := goose.New()
	article, err := g.ExtractFromURL(inputUrl)
	if err != nil {
		return err
	}
	topImageUrl := article.TopImage
	err = DownloadImage(topImageUrl, outputFile)
	return err
}

func DownloadImage(url string, outputFile string) error {
	filepath := outputFile

	out, err := os.Create(filepath)
	defer out.Close()
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
