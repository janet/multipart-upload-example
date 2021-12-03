package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	BUCKET   = "bucket-name"
	KEY      = "janet"
	FILENAME = "lemonade.m4v"
)

func main() {
	// open file
	file, _ := os.Open(FILENAME)
	defer file.Close()

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	s3Client := s3.NewFromConfig(cfg)

	// Create an uploader passing it the client
	uploader := manager.NewUploader(s3Client)

	// Upload
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(KEY),
		Body:   file,
	})
	if err != nil {
		log.Fatal(err)
	}
}
