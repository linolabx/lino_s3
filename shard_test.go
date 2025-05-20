package lino_s3_test

import (
	"testing"

	"github.com/linolabx/lino_s3"
)

func TestHash(t *testing.T) {
	if lino_s3.Hash("mykey") != "9adbe0b3033881f88ebd825bcf763b43" {
		t.Fatal("hash failed")
	}
}

func TestShard(t *testing.T) {
	tmpl2key := map[string]string{
		"{.}":               "mykey",
		"{hash}":            "9adbe0b3033881f88ebd825bcf763b43",
		"{shard.l3}/{.}":    "9a/db/e0/mykey",
		"{shard}/{.}":       "9a/db/e0/mykey",
		"{shard.l4}/{.}":    "9a/db/e0/b3/mykey",
		"{shard.l2}/{hash}": "9a/db/9adbe0b3033881f88ebd825bcf763b43",
	}

	for tmpl, key := range tmpl2key {
		result := lino_s3.ShardT("mykey", tmpl)
		if result != key {
			t.Fatalf("shard template failed: %s, expected: %s, got: %s", tmpl, key, result)
		}
	}

	if lino_s3.ShardT("mykey", "{shard}/%d.%s", 1, "jpg") != "9a/db/e0/1.jpg" {
		t.Fatal("shard template with sprintf failed")
	}
}
