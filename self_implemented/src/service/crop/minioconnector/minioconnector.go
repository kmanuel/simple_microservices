package minioconnector

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go"
	"log"
	"os"
)

const minioHost = "localhost:9000"
const useSsl = false

var accessKey string
var secretKey string
var bucketName string

func Init() {
	godotenv.Load()
	accessKey = os.Getenv("MINIO_ACCESS_KEY")
	secretKey = os.Getenv("MINIO_SECRET_KEY")
	bucketName = os.Getenv("BUCKET_NAME")

}

func DownloadFile(objectName string) string {
	outputFilePath := "/tmp/downloaded" + uuid.New().String() + ".jpg"

	client, err := minio.New(
		minioHost,
		accessKey,
		secretKey,
		useSsl)

	if err != nil {
		log.Fatalln(err)
	}

	err = client.FGetObject(bucketName, objectName, outputFilePath, minio.GetObjectOptions{})

	if err != nil {
		fmt.Println(err)
	}

	return outputFilePath
}

func UploadFile(filePath string) {
	client, err := minio.New(
		minioHost,
		accessKey,
		secretKey,
		useSsl)

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", client)

	bucketExists, err := client.BucketExists(bucketName)

	if err != nil {
		log.Fatalln(err)
	}

	if !bucketExists {
		createBucket(err, client, bucketName)
		log.Printf("successfully created bucket asiatrip")
	} else {
		log.Printf("bucket already exists")
	}

	objectName := uuid.New().String()
	contentType := "img/jpeg"

	n, err := client.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
}

func createBucket(err error, client *minio.Client, bucketName string) {
	err = client.MakeBucket(bucketName, "us-east-1")
	if err != nil {
		log.Fatalln("couldn't create bucket")
	}
}
