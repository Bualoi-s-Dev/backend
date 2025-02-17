package repositories

import (
	"errors"
	"fmt"
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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	genKey := key + "_" + primitive.NewObjectID().Hex()
	_, uploadErr := s.Uploaders.LimitedImgUploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:             aws.String(s.BucketName),
		Key:                aws.String(genKey),
		Body:               strings.NewReader(string(fileBytes)),
		ACL:                types.ObjectCannedACLPublicRead, // Ensure public access
		ContentDisposition: aws.String("inline"),            // Make file viewable in browser
		ContentType:        aws.String(contentType),         // Preserve file type
	})
	if uploadErr != nil {
		log.Println("Error while uploading")
		return "", uploadErr
	}

	return "/" + genKey, nil
}

func (s *S3Repository) DeleteObject(key string) error {
	bucket := s.BucketName
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err := s.Client.DeleteObject(context.TODO(), input)
	fmt.Println("err :", err)
	if err != nil {
		var apiErr *smithy.GenericAPIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case "AccessDenied":
				return fmt.Errorf("access denied: cannot delete object %s from %s", key, bucket)
			default:
				return fmt.Errorf("failed to delete object %s from %s: %v", key, bucket, err)
			}
		}
		return fmt.Errorf("unexpected error deleting object %s from %s: %v", key, bucket, err)
	}

	fmt.Printf("Object %s deleted (or did not exist) from %s.\n", key, bucket)
	return err
}
