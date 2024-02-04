package lino_s3_test

import (
	"testing"

	"github.com/linolabx/lino_s3"
	"github.com/linolabx/lino_s3/internal/test"
)

func TestBasicStructureFunction(t *testing.T) {
	sess := test.GetS3Session()
	s3 := lino_s3.NewLinoS3(sess)
	bucket := s3.Bucket("lino-stor")
	path := bucket.SubPath("text:v1")
	object := path.Object("test.txt")
	object.WriteString("Hello, World!")

	if str, err := object.ReadString(); str != "Hello, World!" {
		t.Fatal(err)
	}

	if _, err := object.Delete(); err != nil {
		t.Fatal(err)
	}
}
