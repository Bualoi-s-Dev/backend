package repositories

import (
	"errors"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type S3Repository struct {
	BucketName string
	Client     *s3.Client
	Uploaders  *S3Uploaders
}
type S3Uploaders struct {
	DefaultUploader    *manager.Uploader
	LimitedImgUploader *manager.Uploader
}

func NewUploaders(client *s3.Client) *S3Uploaders {
	DefaultUploader := manager.NewUploader(client)
	LimitedImgUploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.PartSize = 8 << 20 // 8 MiB
	})
	return &S3Uploaders{
		DefaultUploader:    DefaultUploader,
		LimitedImgUploader: LimitedImgUploader,
	}
}

func NewS3Repository() *S3Repository {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatalln("error:", err)
	}
	bucketName := os.Getenv("S3_BUCKET_NAME")
	// client := s3.NewFromConfig(cfg)
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("BUCKET_URL"))
	})
	uploaders := NewUploaders(client)

	return &S3Repository{
		BucketName: bucketName,
		Client:     client,
		Uploaders:  uploaders,
	}
}

func (s *S3Repository) UploadFile(file *multipart.FileHeader, key string) (string, error) {
	uploadFile, err := file.Open()
	defer uploadFile.Close()

	if err != nil {
		log.Println("Error while opening the file.")
		return "", err
	}

	_, uploadErr := s.Uploaders.LimitedImgUploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:             aws.String(s.BucketName),
		Key:                aws.String(key),
		Body:               uploadFile,
		ACL:                types.ObjectCannedACLPublicRead,             // Ensure public access
		ContentDisposition: aws.String("inline"),                        // Make file viewable in browser
		ContentType:        aws.String(file.Header.Get("Content-Type")), // Preserve file type
	})

	if uploadErr != nil {
		log.Println("Error while uploading")
		return "", uploadErr
	}

	return key, nil
}

func (s *S3Repository) UploadBase64(fileBytes []byte, key string, contentType string) (string, error) {
	_, uploadErr := s.Uploaders.LimitedImgUploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:             aws.String(s.BucketName),
		Key:                aws.String(key),
		Body:               strings.NewReader(string(fileBytes)),
		ACL:                types.ObjectCannedACLPublicRead, // Ensure public access
		ContentDisposition: aws.String("inline"),            // Make file viewable in browser
		ContentType:        aws.String(contentType),         // Preserve file type
	})
	if uploadErr != nil {
		log.Println("Error while uploading")
		return "", uploadErr
	}

	return key, nil
}

func (s *S3Repository) DeleteObject(key string) error {
	bucket := s.BucketName
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err := s.Client.DeleteObject(context.TODO(), input)
	if err != nil {
		var noKey *types.NoSuchKey
		var apiErr *smithy.GenericAPIError
		if errors.As(err, &noKey) {
			log.Printf("Object %s does not exist in %s.\n", key, bucket)
			err = noKey
		} else if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case "AccessDenied":
				log.Printf("Access denied: cannot delete object %s from %s.\n", key, bucket)
				err = nil
				// case "InvalidArgument":
				// 	if bypassGovernance {
				// 		log.Printf("You cannot specify bypass governance on a bucket without lock enabled.")
				// 		err = nil
				// 	}
			}
		}
	}
	return err
}
