package lino_s3_map_test

import (
	"crypto/sha1"
	"fmt"
	"sync"
	"testing"

	"github.com/linolabx/lino_s3/internal/test"
	"github.com/linolabx/lino_s3/structures/lino_s3_map"
)

type TestPiece struct {
	data []byte
	key  string
}

func hash(data []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(data))
}

func TestMap(t *testing.T) {
	bucket := test.GetS3Bucket()
	obj := test.GetS3Object("map-v1", "test.bin")
	tMap := lino_s3_map.NewMap(obj)
	t.Cleanup(func() { tMap.Delete() })

	itemsNum := test.RandomInRange(11, 20)
	tPieces := []TestPiece{}
	for i := 0; i < int(itemsNum); i++ {
		data := test.RandomTestBytes()
		tPieces = append(tPieces, TestPiece{data, hash(data)})
	}

	var wg = new(sync.WaitGroup)

	for _, tPiece := range tPieces {
		wg.Add(1)
		go func(tPiece TestPiece) {
			defer wg.Done()
			if err := tMap.Set(tPiece.key, tPiece.data); err != nil {
				t.Fatalf("set failed: %v", err)
			}
		}(tPiece)
	}
	wg.Wait()

	error := tMap.Save()
	if error != nil {
		t.Fatal(error)
	}

	// test set frozen map
	if err := tMap.Set("test", []byte("test")); err == nil {
		t.Fatal("expected error, got nil")
	}

	partMap, err := lino_s3_map.LoadMap(obj)
	if err != nil {
		t.Fatalf("load map failed: %v", err)
	}

	entireMap, err := lino_s3_map.LoadEntireMap(obj)
	if err != nil {
		t.Fatalf("load map failed: %v", err)
	}

	for _, tPiece := range tPieces {
		mapKey := tPiece.key
		pieceKey, err := tMap.PieceKey(mapKey)
		if err != nil {
			t.Fatalf("get piece key failed: %v", err)
		}

		// read from original lino map
		data, err := tMap.Get(tPiece.key)
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}

		test.BytesCompare(t, tPiece.data, data)

		// read directly from s3 piece
		piece, err := bucket.PieceByKey(pieceKey)
		if err != nil {
			t.Fatalf("get piece failed: %v", err)
		}
		data2, err := piece.ReadBuffer()
		if err != nil {
			t.Fatalf("read piece failed: %v", err)
		}
		test.BytesCompare(t, tPiece.data, data2)

		// read from new lino map instances
		data3, err := partMap.Get(mapKey)
		if err != nil {
			t.Fatalf("get part map failed: %v", err)
		}
		test.BytesCompare(t, tPiece.data, data3)

		data4, err := entireMap.Get(mapKey)
		if err != nil {
			t.Fatalf("get full map failed: %v", err)
		}
		test.BytesCompare(t, tPiece.data, data4)
	}
}
