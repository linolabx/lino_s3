package test

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/geektheripper/vast-dsn/s3_dsn"
	"github.com/linolabx/lino_s3"
)

func GetS3Session() *session.Session {
	config, err := s3_dsn.ParseS3DSN("s3://minioadmin:minioadmin@localhost:9000?s3-force-path-style=true&protocol=http")
	if err != nil {
		panic(err)
	}

	sess, err := session.NewSession(config)
	if err != nil {
		panic(err)
	}

	return sess
}

func GetS3Bucket() *lino_s3.LinoS3Bucket {
	return lino_s3.NewLinoS3(GetS3Session()).Bucket("lino-stor")
}

func GetS3Object(subPath string, key string) *lino_s3.LinoS3Object {
	return GetS3Bucket().SubPath(subPath).Object(key)
}
