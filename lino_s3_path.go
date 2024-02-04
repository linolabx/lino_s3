package lino_s3

import (
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws/session"
)

type LinoS3Path struct {
	sess   *session.Session
	bucket string
	path   string

	interceptors []*Interceptor
}

func (s *LinoS3Path) UseInterceptors(interceptors ...*Interceptor) *LinoS3Path {
	s.interceptors = interceptors
	return s
}

func (s *LinoS3Path) SubPath(path string) *LinoS3Path {
	return &LinoS3Path{
		sess:         s.sess,
		bucket:       s.bucket,
		path:         filepath.Join(s.path, path),
		interceptors: s.interceptors,
	}
}

func (s *LinoS3Path) Object(key string) *LinoS3Object {
	return &LinoS3Object{
		sess:         s.sess,
		bucket:       s.bucket,
		key:          filepath.Join(s.path, key),
		interceptors: s.interceptors,
	}
}
