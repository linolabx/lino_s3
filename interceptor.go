package lino_s3

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type transformer[T any] func(input T) (T, error)

type Interceptor struct {
	PreHead    transformer[*s3.HeadObjectInput]
	PostHead   transformer[*s3.HeadObjectOutput]
	PreGet     transformer[*s3.GetObjectInput]
	PostGet    transformer[*s3.GetObjectOutput]
	PrePut     transformer[*s3.PutObjectInput]
	PostPut    transformer[*s3.PutObjectOutput]
	PreDelete  transformer[*s3.DeleteObjectInput]
	PostDelete transformer[*s3.DeleteObjectOutput]
	PreUpload  transformer[*s3manager.UploadInput]
	PostUpload transformer[*s3manager.UploadOutput]
}

type callInterceptor[T any] func(intcpt *Interceptor, v T) (T, error)

// Reverse Resolve
func preResolve[T interface{}](intcpts []*Interceptor, call callInterceptor[T], v T) (T, error) {
	for i := len(intcpts) - 1; i >= 0; i-- {
		var err error
		v, err = call(intcpts[i], v)
		if err != nil {
			return v, err
		}
	}
	return v, nil
}

func postResolve[T interface{}](intcpts []*Interceptor, call callInterceptor[T], v T) (T, error) {
	for _, intcpt := range intcpts {
		var err error
		v, err = call(intcpt, v)
		if err != nil {
			return v, err
		}
	}
	return v, nil
}

func usePreGet(intcpt *Interceptor, value *s3.GetObjectInput) (*s3.GetObjectInput, error) {
	if intcpt.PreGet == nil {
		return value, nil
	}
	return intcpt.PreGet(value)
}

func usePostGet(intcpt *Interceptor, value *s3.GetObjectOutput) (*s3.GetObjectOutput, error) {
	if intcpt.PostGet == nil {
		return value, nil
	}
	return intcpt.PostGet(value)
}

func usePreHead(intcpt *Interceptor, value *s3.HeadObjectInput) (*s3.HeadObjectInput, error) {
	if intcpt.PreHead == nil {
		return value, nil
	}
	return intcpt.PreHead(value)
}

func usePostHead(intcpt *Interceptor, value *s3.HeadObjectOutput) (*s3.HeadObjectOutput, error) {
	if intcpt.PostHead == nil {
		return value, nil
	}
	return intcpt.PostHead(value)
}

func usePrePut(intcpt *Interceptor, value *s3.PutObjectInput) (*s3.PutObjectInput, error) {
	if intcpt.PrePut == nil {
		return value, nil
	}
	return intcpt.PrePut(value)
}

func usePostPut(intcpt *Interceptor, value *s3.PutObjectOutput) (*s3.PutObjectOutput, error) {
	if intcpt.PostPut == nil {
		return value, nil
	}
	return intcpt.PostPut(value)
}

func usePreDelete(intcpt *Interceptor, value *s3.DeleteObjectInput) (*s3.DeleteObjectInput, error) {
	if intcpt.PreDelete == nil {
		return value, nil
	}
	return intcpt.PreDelete(value)
}

func usePostDelete(intcpt *Interceptor, value *s3.DeleteObjectOutput) (*s3.DeleteObjectOutput, error) {
	if intcpt.PostDelete == nil {
		return value, nil
	}
	return intcpt.PostDelete(value)
}

func usePreUpload(intcpt *Interceptor, value *s3manager.UploadInput) (*s3manager.UploadInput, error) {
	if intcpt.PreUpload == nil {
		return value, nil
	}
	return intcpt.PreUpload(value)
}

func usePostUpload(intcpt *Interceptor, value *s3manager.UploadOutput) (*s3manager.UploadOutput, error) {
	if intcpt.PostUpload == nil {
		return value, nil
	}
	return intcpt.PostUpload(value)
}
