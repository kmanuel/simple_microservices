package function

import (
	"bytes"
	"github.com/advancedlogic/GoOse"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/prometheus/common/log"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const outputImageLocation = "/tmp/"

type Task struct {
	ID  string `jsonapi:"primary,most_significant_image_task"`
	Url string `jsonapi:"attr,url"`
}

func Handle(req []byte) string {

	minioService := initMinio()

	task := new(Task)
	task.ID = uuid.New().String()

	err := jsonapi.UnmarshalPayload(bytes.NewReader(req), task)
	if err != nil {
		panic(err)
	}

	err = handleTask(task, *minioService)
	if err != nil {
		panic(err)
	}

	return ""
}

func initMinio() *minioconnector.MinioService {
	minioHost := os.Getenv("MINIO_HOST")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecret := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("INPUT_BUCKET_NAME")

	log.Errorf("initializing minio with host=%s accessKey=%s secret=%s bucketName=%s", minioHost, minioAccessKey, minioSecret, bucketName)
	return minioconnector.NewMinioService(
		minioHost,
		minioAccessKey,
		minioSecret,
		bucketName,
		"mostsignificantimage")
}

func handleTask(t *Task, minioService minioconnector.MinioService) error {
	task := t

	outputFile := outputImageLocation + uuid.New().String() + ".jpg"

	err := ExtractMostSignificantImage(task.Url, outputFile)
	if err != nil {
		return err
	}

	_, err = minioService.UploadFileWithName(outputFile, createFileName(t))
	if err != nil {
		return err
	}

	return nil
}

func createFileName(task *Task) string {
	inputFileName := strings.Replace(task.Url, "http://", "", -1)
	inputFileName = strings.Replace(inputFileName, "https://", "", -1)
	inputFileName = strings.Replace(inputFileName, ".", "_", -1)
	inputFileName = strings.Replace(inputFileName, "/", "_", -1)
	timestamp := "_" + strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	return inputFileName + timestamp + ".jpg"
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
