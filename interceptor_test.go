package lino_s3_test

import (
	"testing"

	"github.com/linolabx/lino_s3/interceptors"
	"github.com/linolabx/lino_s3/internal/test"
)

func TestMixedInterceptors(t *testing.T) {
	book, err := test.GetBook()
	if err != nil {
		t.Fatal(err)
	}

	bookObj := test.
		GetS3Object("cbor-v2", "book.cbor.zstd.gz").
		UseInterceptors(
			interceptors.Gzip,
			interceptors.Zstd,
		)

	t.Cleanup(func() { bookObj.Delete() })

	if err := bookObj.WriteCBOR(book); err != nil {
		t.Fatal(err)
	}

	_book := test.Book{}
	if err := bookObj.ReadCBOR(&_book); err != nil {
		t.Fatal(err)
	}

	if _book.Content != book.Content {
		t.Fatal("Book CBOR read failed")
	}
}
