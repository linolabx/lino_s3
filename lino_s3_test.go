package lino_s3_test

import (
	"testing"

	"github.com/linolabx/lino_s3"
	"github.com/linolabx/lino_s3/internal/test"
)

func TestBasicStructureFunction(t *testing.T) {
	client := test.GetS3Client()
	s3 := lino_s3.NewLinoS3(client)
	bucket := s3.Bucket("lino-stor")
	path := bucket.SubPath("text-v1")
	object := path.Object("test.txt")
	if err := object.WriteString("Hello, World!"); err != nil {
		t.Fatal(err)
	}

	if str, err := object.ReadString(); str != "Hello, World!" {
		t.Fatal(err)
	}

	if _, err := object.Delete(); err != nil {
		t.Fatal(err)
	}
}
