package service

import (
	"github.com/advancedlogic/GoOse"
	"github.com/google/uuid"
	"github.com/kmanuel/minioconnector"
	"github.com/kmanuel/simple_microservices/go_kit/service/most_significant_image/src/model"
	"github.com/prometheus/common/log"
	"io"
	"net/http"
	"os"
)

const outputImageLocation = "/tmp/"

type ImageService interface {
	HandleTask(*model.Task) error
}

type mostSignificantImageService struct{
	minioService minioconnector.MinioService
}

func NewMostSignificantImageService(service *minioconnector.MinioService) ImageService {
	return mostSignificantImageService{*service}
}

func (s mostSignificantImageService) HandleTask(t *model.Task) error {
	g := goose.New()
	article, err := g.ExtractFromURL(t.Url)
	if err != nil {
		return err
	}
	filePath := outputImageLocation + uuid.New().String()
	topImageUrl := article.TopImage
	if err = downloadImage(topImageUrl, filePath); err != nil {
		return err
	}

	log.Info("uploading file")
	if _, err = s.minioService.UploadFileWithName(filePath, t.ID); err != nil {
		log.Info("error while uploading file", err)
		return err
	}

	return nil
}

func downloadImage(url string, outputFile string) error {
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
