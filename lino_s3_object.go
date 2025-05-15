package lino_s3

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fxamacker/cbor/v2"
	"github.com/gocarina/gocsv"
	"github.com/linolabx/lino_s3/internal"
)

type LinoS3Object struct {
	client *s3.Client
	bucket string
	key    string

	interceptors []*Interceptor
}

func (s *LinoS3Object) Piece(start int64, end int64) *LinoS3Piece {
	return &LinoS3Piece{s.client, s.bucket, s.key, start, end}
}

func (s *LinoS3Object) UseInterceptors(interceptors ...*Interceptor) *LinoS3Object {
	s.interceptors = interceptors
	return s
}

func (s *LinoS3Object) HasInterceptors() bool {
	return len(s.interceptors) > 0
}

func (s *LinoS3Object) Get() (*s3.GetObjectOutput, error) {
	input, err := preResolve(s.interceptors, usePreGet,
		&s3.GetObjectInput{Bucket: &s.bucket, Key: &s.key},
	)

	if err != nil {
		return nil, err
	}

	out, err := s.client.GetObject(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return postResolve(s.interceptors, usePostGet, out)
}

func (s *LinoS3Object) PresignGet(optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	presigner := s3.NewPresignClient(s.client)

	return presigner.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &s.key,
	}, optFns...)
}

func (s *LinoS3Object) Head() (*s3.HeadObjectOutput, error) {
	input, err := preResolve(s.interceptors, usePreHead,
		&s3.HeadObjectInput{Bucket: &s.bucket, Key: &s.key},
	)
	if err != nil {
		return nil, err
	}

	out, err := s.client.HeadObject(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return postResolve(s.interceptors, usePostHead, out)
}

func (s *LinoS3Object) Put(input s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	input.Bucket = &s.bucket
	input.Key = &s.key

	_input, err := preResolve(s.interceptors, usePrePut, &input)
	if err != nil {
		return nil, err
	}

	out, err := s.client.PutObject(context.Background(), _input)
	if err != nil {
		return nil, err
	}

	return postResolve(s.interceptors, usePostPut, out)
}

func (s *LinoS3Object) PresignPut(input s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	presigner := s3.NewPresignClient(s.client)

	input.Bucket = &s.bucket
	input.Key = &s.key

	return presigner.PresignPutObject(context.Background(), &input, optFns...)
}

func (s *LinoS3Object) Upload(input s3.PutObjectInput) (*manager.UploadOutput, error) {
	input.Bucket = &s.bucket
	input.Key = &s.key

	_input, err := preResolve(s.interceptors, usePreUpload, &input)
	if err != nil {
		return nil, err
	}

	uploader := manager.NewUploader(s.client, func(u *manager.Uploader) {
		u.PartSize = 10 << 20
		u.Concurrency = 5
	})

	out, err := uploader.Upload(context.Background(), _input)
	if err != nil {
		return nil, err
	}

	return postResolve(s.interceptors, usePostUpload, out)
}

func (s *LinoS3Object) Delete() (*s3.DeleteObjectOutput, error) {
	input, err := preResolve(s.interceptors, usePreDelete,
		&s3.DeleteObjectInput{Bucket: &s.bucket, Key: &s.key},
	)
	if err != nil {
		return nil, err
	}

	out, err := s.client.DeleteObject(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return postResolve(s.interceptors, usePostDelete, out)
}

func (s *LinoS3Object) ReadTo(writer io.WriteCloser) error {
	out, err := s.Get()
	if err != nil {
		return err
	}

	defer out.Body.Close()
	defer writer.Close()

	_, err = io.Copy(writer, out.Body)
	return err
}

func (s *LinoS3Object) ReadBuffer() ([]byte, error) {
	out, err := s.Get()
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()

	return io.ReadAll(out.Body)
}

func (s *LinoS3Object) ReadString() (string, error) {
	data, err := s.ReadBuffer()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *LinoS3Object) ReadJSON(v interface{}) error {
	data, err := s.ReadBuffer()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func (s *LinoS3Object) ReadCBOR(v interface{}) error {
	data, err := s.ReadBuffer()
	if err != nil {
		return err
	}

	return cbor.Unmarshal(data, v)
}

func (s *LinoS3Object) ReadCSV(v interface{}) error {
	out, err := s.Get()
	if err != nil {
		return err
	}

	defer out.Body.Close()

	return gocsv.Unmarshal(out.Body, v)
}

func (s *LinoS3Object) WriteFrom(reader io.ReadCloser, contentType ...string) error {
	_, err := s.Upload(s3.PutObjectInput{
		Body:        reader,
		ContentType: internal.OptionalPointer(contentType...),
	})
	return err
}

func (s *LinoS3Object) WriteBuffer(data []byte, contentType ...string) error {
	_, err := s.Put(s3.PutObjectInput{
		Body:        bytes.NewReader(data),
		ContentType: internal.OptionalPointer(contentType...),
	})
	return err
}

func (s *LinoS3Object) WriteString(_string string, contentType ...string) error {
	return s.WriteBuffer([]byte(_string), contentType...)
}

func (s *LinoS3Object) WriteJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return s.WriteBuffer(data, "application/json")
}

func (s *LinoS3Object) WriteCBOR(v interface{}) error {
	data, err := cbor.Marshal(v)
	if err != nil {
		return err
	}

	return s.WriteBuffer(data, "application/cbor")
}

func (s *LinoS3Object) WriteCSV(v interface{}) error {
	data, err := gocsv.MarshalBytes(v)
	if err != nil {
		return err
	}

	return s.WriteBuffer(data, "text/csv")
}

func (s *LinoS3Object) Key() string {
	return s.key
}
