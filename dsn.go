package lino_s3

import (
	"github.com/geektheripper/vast-dsn/s3_client"
	"github.com/geektheripper/vast-dsn/s3_dsn"
	"github.com/geektheripper/vast-dsn/utils"
)

func LoadS3(dsn string) (*LinoS3, error) {
	opts, err := s3_dsn.ParseS3(dsn)
	if err != nil {
		return nil, err
	}

	client, err := s3_client.NewS3Client(opts)
	if err != nil {
		return nil, err
	}

	return NewLinoS3(client), nil
}

func LoadS3Bucket(dsn string) (*LinoS3Bucket, error) {
	opts, bucket, err := s3_dsn.ParseS3Bucket(dsn)
	if err != nil {
		return nil, err
	}

	client, err := s3_client.NewS3Client(opts)
	if err != nil {
		return nil, err
	}

	return NewLinoS3(client).Bucket(bucket), nil
}

func LoadS3Object(dsn string) (*LinoS3Object, error) {
	opts, bucket, key, err := s3_dsn.ParseS3Object(dsn)
	if err != nil {
		return nil, err
	}

	client, err := s3_client.NewS3Client(opts)
	if err != nil {
		return nil, err
	}

	return NewLinoS3(client).Bucket(bucket).Object(key), nil
}

func LoadS3Path(dsn string) (*LinoS3Path, error) {
	opts, bucket, path, err := s3_dsn.ParseS3Path(dsn)
	if err != nil {
		return nil, err
	}

	client, err := s3_client.NewS3Client(opts)
	if err != nil {
		return nil, err
	}

	return NewLinoS3(client).Bucket(bucket).SubPath(path), nil
}
func MustLoadS3(dsn string, logger ...utils.Logger) *LinoS3 {
	log := utils.EnsureLogger(logger...)

	client, err := LoadS3(dsn)
	if err != nil {
		log.Fatalf("failed to load s3: %v", err)
	}

	return client
}

func MustLoadS3Bucket(dsn string, logger ...utils.Logger) *LinoS3Bucket {
	log := utils.EnsureLogger(logger...)

	client, err := LoadS3Bucket(dsn)
	if err != nil {
		log.Fatalf("failed to load s3 bucket: %v", err)
	}

	return client
}

func MustLoadS3Object(dsn string, logger ...utils.Logger) *LinoS3Object {
	log := utils.EnsureLogger(logger...)

	client, err := LoadS3Object(dsn)
	if err != nil {
		log.Fatalf("failed to load s3 object: %v", err)
	}
	return client
}

func MustLoadS3Path(dsn string, logger ...utils.Logger) *LinoS3Path {
	log := utils.EnsureLogger(logger...)

	client, err := LoadS3Path(dsn)
	if err != nil {
		log.Fatalf("failed to load s3 path: %v", err)
	}
	return client
}
