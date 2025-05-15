package lino_s3_test

import (
	"bytes"
	"testing"

	"github.com/linolabx/lino_s3/internal/test"
)

type TestPiece struct {
	start int64
	end   int64
	data  []byte
}

func TestS3Piece(t *testing.T) {
	obj := test.GetS3Object("piece-v1", "test.bin")
	t.Cleanup(func() { obj.Delete() })

	itemsNum := test.RandomInRange(11, 20)
	tPieces := []TestPiece{}
	objData := []byte{}
	offset := 0
	for i := 0; i < int(itemsNum); i++ {
		data := test.RandomTestBytes()
		tPieces = append(tPieces, TestPiece{int64(offset), int64(offset + len(data)), data})
		objData = append(objData, data...)
		offset += len(data)
	}

	obj.WriteBuffer(objData)

	for _, tPiece := range tPieces {
		data, err := obj.Piece(tPiece.start, tPiece.end).ReadBuffer()
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(data, tPiece.data) {
			t.Fatalf(
				"piece data not equal, expected %v...%v(%d), got %v...%v(%d)",
				tPiece.data[:16], tPiece.data[len(tPiece.data)-16:], len(tPiece.data),
				data[:16], data[len(data)-16:], len(data),
			)
		}
	}

	tPiece := tPieces[10]
	pieceKey := obj.Piece(tPiece.start, tPiece.end).Key()
	t.Logf("piece key: %s", pieceKey)

	s3Piece, err := test.GetS3Bucket().PieceByKey(pieceKey)
	if err != nil {
		t.Fatal(err)
	}

	data, err := s3Piece.ReadBuffer()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, tPiece.data) {
		t.Fatalf(
			"piece data not equal, expected %v...%v(%d), got %v...%v(%d)",
			tPiece.data[:16], tPiece.data[len(tPiece.data)-16:], len(tPiece.data),
			data[:16], data[len(data)-16:], len(data),
		)
	}
}
