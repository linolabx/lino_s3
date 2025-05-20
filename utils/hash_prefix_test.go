package utils_test

import (
	"testing"

	"github.com/linolabx/lino_s3/utils"
)

func TestHashSplit(t *testing.T) {
	if utils.HashSplit("123456") != "e1/0a/dc" ||
		utils.HashSplit("123456", 5) != "e1/0a/dc/39/49" ||
		utils.HashSplit("123456", 16) != "e1/0a/dc/39/49/ba/59/ab/be/56/e0/57/f2/0f/88/3e" {
		t.Fatal("hash split failed")
	}

	// HashTemplate("mykey", "{.}") => "mykey"
	// HashTemplate("mykey", "{hash}") => "9adbe0b3033881f88ebd825bcf763b43"
	// HashTemplate("mykey", "{hash.l3}/{.}") => "9a/db/e0/mykey"
	// HashTemplate("mykey", "{hash.l4}/{.}") => "9a/db/e0/b3/mykey"
	// HashTemplate("mykey", "{hash.l3}/{hash}") => "9a/db/e0/9adbe0b3033881f88ebd825bcf763b43"

	if utils.HashTemplate("mykey", "{.}") != "mykey" ||
		utils.HashTemplate("mykey", "{hash}") != "9adbe0b3033881f88ebd825bcf763b43" ||
		utils.HashTemplate("mykey", "{hash.l3}/{.}") != "9a/db/e0/mykey" ||
		utils.HashTemplate("mykey", "{hash.l4}/{.}") != "9a/db/e0/b3/mykey" ||
		utils.HashTemplate("mykey", "{hash.l3}/{hash}") != "9a/db/e0/9adbe0b3033881f88ebd825bcf763b43" {
		t.Fatal("hash template failed")
	}
}

func TestHashPrefix(t *testing.T) {
	if utils.HashPrefix("hello.txt", 3) != "2e/54/14/hello.txt" {
		t.Fatal("hash prefix failed")
	}
}
