package lino_s3

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type LinoS3 struct {
	client *s3.Client

	interceptors []*Interceptor
}

func (s *LinoS3) UseInterceptors(interceptors ...*Interceptor) *LinoS3 {
	s.interceptors = interceptors
	return s
}

func NewLinoS3(client *s3.Client) *LinoS3 {
	return &LinoS3{client: client}
}

func (s *LinoS3) Bucket(bucketname string) *LinoS3Bucket {
	return &LinoS3Bucket{s.client, bucketname, s.interceptors}
}
