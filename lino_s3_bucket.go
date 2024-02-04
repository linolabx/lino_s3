package lino_s3

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

type LinoS3Bucket struct {
	sess   *session.Session
	bucket string

	interceptors []*Interceptor
}

func (s *LinoS3Bucket) UseInterceptors(interceptors ...*Interceptor) *LinoS3Bucket {
	s.interceptors = interceptors
	return s
}

func (s *LinoS3Bucket) SubPath(path string) *LinoS3Path {
	return &LinoS3Path{
		sess:         s.sess,
		bucket:       s.bucket,
		path:         path,
		interceptors: s.interceptors,
	}
}

func (s *LinoS3Bucket) Object(key string) *LinoS3Object {
	return &LinoS3Object{
		sess:         s.sess,
		bucket:       s.bucket,
		key:          key,
		interceptors: s.interceptors,
	}
}
