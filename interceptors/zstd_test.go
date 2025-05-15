package interceptors_test

import (
	"testing"

	"github.com/linolabx/lino_s3/interceptors"
	"github.com/linolabx/lino_s3/internal/test"
)

func TestZstd(t *testing.T) {
	book, err := test.GetBook()
	if err != nil {
		t.Fatal(err)
	}

	bookObj := test.
		GetS3Object("cbor-v2", "book.cbor.zstd").
		UseInterceptors(interceptors.Zstd)

	t.Cleanup(func() { bookObj.Delete() })

	if err := bookObj.WriteCBOR(book); err != nil {
		t.Fatal(err)
	}

	if resp, err := bookObj.Head(); err != nil {
		t.Fatal(err)
	} else if *resp.ContentLength >= (book.Size / 2) {
		t.Fatal("compression failed")
	} else {
		// t.Logf("compress rate: %d%%", *resp.ContentLength/book.Size*100)
	}

	var _book test.Book
	if err := bookObj.ReadCBOR(&_book); err != nil {
		t.Fatal(err)
	}

	if _book.Content != book.Content {
		t.Fatal("Book CBOR read failed")
	}
}
