package minioconnector

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

const useSsl = false

var minioHost string
var accessKey string
var secretKey string
var bucketName string

func Init(
	minioHostArg string,
	accessKeyArg string,
	secretKeyArg string,
	bucketNameArg string) {

	minioHost = minioHostArg
	accessKey = accessKeyArg
	secretKey = secretKeyArg
	bucketName = bucketNameArg

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

func UploadFile(filePath string) string {
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
		log.WithField("bucketname", bucketName).Info("successfully created bucket")
	} else {
		log.WithField("bucketname", bucketName).Debug("bucket already exists")
	}

	objectName := uuid.New().String()
	contentType := "img/jpeg"

	n, err := client.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.WithField("objectName", objectName).Info("Successfully uploaded %s of size %d\n", objectName, n)

	return objectName
}

func createBucket(err error, client *minio.Client, bucketName string) {
	err = client.MakeBucket(bucketName, "us-east-1")
	if err != nil {
		log.WithField("bucketName", bucketName).Fatalln("couldn't create bucket")
	}
}
