package s3

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"net/http"
)

type s3Uploader struct {
	awsSession *session.Session
	bucketName string
}

func Build(S3Region string, bucketName string) (*s3Uploader, error) {
	s := s3Uploader{bucketName: bucketName}
	s3Session, err := session.NewSession(&aws.Config{Region: aws.String(S3Region),
		Credentials: credentials.NewEnvCredentials()},
	)
	if err != nil {
		return nil, fmt.Errorf("there was an error creating a s3 session: %w", err)
	}
	s.awsSession = s3Session
	return &s, nil
}

func (s s3Uploader) AddFileToS3(buffer []byte, fileName string) error {
	_, err := s3.New(s.awsSession).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(s.bucketName),
		Key:                aws.String(fileName),
		ACL:                aws.String("private"),
		Body:               bytes.NewReader(buffer),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
	})
	return err
}
