package services

import (
	"log"
	"mime/multipart"
	"os"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	BucketName string
	Client     *s3.Client
}

func NewS3Service() *S3Service {
	cfg, err := config.LoadDefaultConfig((context.TODO()))
	if err != nil {
		log.Fatalln("error:", err)
	}
	bucketName := os.Getenv("S3_BUCKET_NAME")
	client := s3.NewFromConfig(cfg)

	return &S3Service{
		BucketName: bucketName,
		Client:     client,
	}
}

func (s *S3Service) UploadFile(file *multipart.FileHeader) (string, error) {
	uploader := manager.NewUploader(s.Client, func(u *manager.Uploader) {
		u.PartSize = 8 << 20 // 8 MiB
	})

	uploadFile, err := file.Open()
	defer uploadFile.Close()

	if err != nil {
		log.Println("Error while opening the file.")
		return "", err
	}

	result, uploadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(file.Filename),
		Body:   uploadFile,
	})
	if uploader != nil {
		log.Println("Error while uploading")
		return "", err
	}

	return result.Location, uploadErr

}

// cfg, err := config.LoadDefaultConfig(context.TODO())
// if err != nil {
//     log.Printf("error: %v", err)
//     return
// }

// client := s3.NewFromConfig(cfg)

// uploader := manager.NewUploader(client)
// result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
//     Bucket: aws.String("amzn-s3-demo-bucket"),
//     Key:    aws.String("my-object-key"),
//     Body:   uploadFile,
// })
