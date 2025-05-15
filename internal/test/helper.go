package test

import (
	"fmt"

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
	return lino_s3.MustLoadS3Bucket("s3://minioadmin:minioadmin@localhost:9000/lino-stor?use-path-style=true&protocol=http")
}

func GetS3Object(subPath string, key string) *lino_s3.LinoS3Object {
	return lino_s3.MustLoadS3Object(fmt.Sprintf("s3://minioadmin:minioadmin@localhost:9000/lino-stor/%s/%s?use-path-style=true&protocol=http", subPath, key))
}
