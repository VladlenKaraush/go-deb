package s3

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	s3Conn *s3.Client
}

type S3AuthCreds struct {
	Key    string
	Secret string
}

func GetS3Client(authCreds S3AuthCreds) S3Client {
	const defaultRegion = "us-east-1"
	staticResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               "http://localhost:9000",
				SigningRegion:     defaultRegion,
				HostnameImmutable: true,
			}, nil
		})
	cfg := aws.Config{
		Region:                      defaultRegion,
		Credentials:                 credentials.NewStaticCredentialsProvider(authCreds.Key, authCreds.Secret, ""),
		EndpointResolverWithOptions: staticResolver,
	}

	s3Conn := s3.NewFromConfig(cfg)
	return S3Client{s3Conn: s3Conn}
}

func (cl *S3Client) UploadPackage(body []byte, bucket, key string) {
	cl.s3Conn.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})
}
