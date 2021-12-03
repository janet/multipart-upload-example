// based off of https://github.com/mathisve/golang-S3-Multipart-Upload-Example/

package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	BUCKET    = "bucket-name"
	KEY       = "janet"
	FILENAME  = "lemonade.m4v"
	PART_SIZE = 6_000_000
	RETRIES   = 2
)

var (
	s3Client *s3.Client
)

func init() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	s3Client = s3.NewFromConfig(cfg)
}

func main() {
	// open file
	file, _ := os.Open(FILENAME)
	defer file.Close()

	// Get file size
	stats, _ := file.Stat()
	fileSize := stats.Size()

	// put file in byteArray
	buffer := make([]byte, fileSize) // wouldn't want to do this for a large file because it would store a potentially super large file into memory
	file.Read(buffer)

	// start multipart upload
	createdResp, err := s3Client.CreateMultipartUpload(context.TODO(), &s3.CreateMultipartUploadInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(KEY),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("")
	log.Printf("Upload ID: %d", createdResp.UploadId)

	var start, currentSize int
	var remaining = int(fileSize)
	var partNum = 1
	var completedParts []s3Types.CompletedPart

	for start = 0; remaining != 0; start += PART_SIZE {
		if remaining < PART_SIZE {
			currentSize = remaining
		} else {
			currentSize = PART_SIZE
		}
		completed, err := Upload(createdResp, buffer[start:start+currentSize], partNum)
		if err != nil {
			_, err = s3Client.AbortMultipartUpload(context.TODO(), &s3.AbortMultipartUploadInput{
				Bucket:   createdResp.Bucket,
				Key:      createdResp.Key,
				UploadId: createdResp.UploadId,
			})
			if err != nil {
				log.Fatal(err)
			}
		}
		remaining -= currentSize
		fmt.Printf("Part %v complete, %v bytes remaining\n", partNum, remaining)
		completedParts = append(completedParts, completed)
		partNum++
	}

	// complete multipart upload
	input := &s3.CompleteMultipartUploadInput{
		Bucket:   createdResp.Bucket,
		Key:      createdResp.Key,
		UploadId: createdResp.UploadId,
		MultipartUpload: &s3Types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}
	result, err := s3Client.CompleteMultipartUpload(context.TODO(), input)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Completed upload")
	log.Println(result)

}

func Upload(resp *s3.CreateMultipartUploadOutput, fileBytes []byte, partNum int) (completedPart s3Types.CompletedPart, err error) {
	var try int
	for try <= RETRIES {
		uploadResp, err := s3Client.UploadPart(context.TODO(), &s3.UploadPartInput{
			Body:          bytes.NewReader(fileBytes),
			Bucket:        aws.String(BUCKET),
			Key:           aws.String(KEY),
			PartNumber:    int32(partNum),
			UploadId:      resp.UploadId,
			ContentLength: int64(len(fileBytes)),
		})
		if err != nil {
			fmt.Println(err)
			if try == RETRIES {
				return s3Types.CompletedPart{}, err
			} else {
				try++
			}
		} else {
			return s3Types.CompletedPart{
				ETag:       uploadResp.ETag,
				PartNumber: int32(partNum),
			}, nil
		}
	}
	return s3Types.CompletedPart{}, nil
}
