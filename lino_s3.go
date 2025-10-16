package lino_s3

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type LinoS3 struct {
	client   *s3.Client
	psClient *s3.PresignClient

	interceptors []*Interceptor
}

func (s *LinoS3) UseInterceptors(interceptors ...*Interceptor) *LinoS3 {
	s.interceptors = interceptors
	return s
}

func (s *LinoS3) UsePresignClient(presignClient *s3.PresignClient) *LinoS3 {
	s.psClient = presignClient
	return s
}

func NewLinoS3(client *s3.Client, presignClient ...*s3.PresignClient) *LinoS3 {
	var psClient *s3.PresignClient
	if len(presignClient) == 0 {
		psClient = s3.NewPresignClient(client)
	} else if len(presignClient) == 1 {
		psClient = presignClient[0]
	} else {
		panic("invalid presign client")
	}

	return &LinoS3{client: client, psClient: psClient}
}

func (s *LinoS3) Bucket(bucketname string) *LinoS3Bucket {
	return &LinoS3Bucket{
		client:       s.client,
		psClient:     s.psClient,
		bucket:       bucketname,
		interceptors: s.interceptors,
	}
}
