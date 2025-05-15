package test

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/geektheripper/vast-dsn/s3_client"
	"github.com/geektheripper/vast-dsn/s3_dsn"
	"github.com/linolabx/lino_s3"
)

func GetS3Client() *s3.Client {
	client, err := s3_client.NewS3Client(s3_dsn.MustParseS3("s3://minioadmin:minioadmin@localhost:9000?use-path-style=true&protocol=http"))
	if err != nil {
		panic(err)
	}

	return client
}

func GetS3Bucket() *lino_s3.LinoS3Bucket {
	return lino_s3.NewLinoS3(GetS3Client()).Bucket("lino-stor")
}

func GetS3Object(subPath string, key string) *lino_s3.LinoS3Object {
	return GetS3Bucket().SubPath(subPath).Object(key)
}
