package lino_s3

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

type LinoS3 struct {
	sess *session.Session

	interceptors []*Interceptor
}

func (s *LinoS3) UseInterceptors(interceptors ...*Interceptor) *LinoS3 {
	s.interceptors = interceptors
	return s
}

func NewLinoS3(sess *session.Session) *LinoS3 {
	return &LinoS3{sess: sess}
}

func (s *LinoS3) Bucket(bucketname string) *LinoS3Bucket {
	return &LinoS3Bucket{s.sess, bucketname, s.interceptors}
}
