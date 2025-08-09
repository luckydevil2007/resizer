package s3storage

import (
	"bytes"
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/luckydevil2007/audionotes/entities"
)

type S3Storage struct {
	Bucket string
	Client *s3.S3
}

func NewS3Storage(Bucket string) *S3Storage {
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentials(
			"minio",
			"minio123",
			"",
		),
		Endpoint:         aws.String("172.18.0.2:9000"),
		Region:           aws.String("RU"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession, err := session.NewSession(&cfg)
	s3Client := s3.New(newSession)
	if err != nil {
		return nil
	}
	return &S3Storage{Bucket: Bucket, Client: s3Client}
}

func (s *S3Storage) Save(ctx context.Context, note *entities.Note) error {
	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Bucket: &[]string{s.Bucket}[0],
		Key:    &note.Path, //сделать хешироние
		Body:   bytes.NewReader(note.Data),
	})

	return err
}

func (s *S3Storage) Delete(ctx context.Context, note *entities.Note) error {
	_, err := s.Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &[]string{s.Bucket}[0],
		Key:    &note.Path,
	})

	return err
}

func (s *S3Storage) Open(ctx context.Context, path string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(path),
	}

	result, err := s.Client.GetObject(input)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}
