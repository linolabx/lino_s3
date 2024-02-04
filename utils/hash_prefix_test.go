package utils_test

import (
	"testing"

	"github.com/linolabx/lino_s3/utils"
)

func ExpectPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	f()
}

func TestHashSplit(t *testing.T) {
	ExpectPanic(t, func() { utils.HashSplit("123456", -1) })
	ExpectPanic(t, func() { utils.HashSplit("123456", 0) })
	ExpectPanic(t, func() { utils.HashSplit("123456", 17) })
	ExpectPanic(t, func() { utils.HashSplit("123456", 2, 3) })

	if utils.HashSplit("123456") != "e1/0a/dc" ||
		utils.HashSplit("123456", 5) != "e1/0a/dc/39/49" ||
		utils.HashSplit("123456", 16) != "e1/0a/dc/39/49/ba/59/ab/be/56/e0/57/f2/0f/88/3e" {
		t.Fatal("hash split failed")
	}
}

func TestHashPrefix(t *testing.T) {
	if utils.HashPrefix("hello.txt", 3) != "2e/54/14/hello.txt" {
		t.Fatal("hash prefix failed")
	}
}
