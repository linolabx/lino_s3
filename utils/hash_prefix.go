package utils

import (
	"crypto/md5"
	"encoding/hex"
	"path/filepath"
	"strings"

	"github.com/linolabx/lino_s3/internal"
)

func HashSplit(s string, level ...int) string {
	_level := internal.OptionalParam(3, level...)
	if _level > 16 || _level <= 0 {
		panic("level too high, md5 hex is 32 characters long, so level should be between 1 and 16 inclusive.")
	}

	h := md5.New()
	h.Write([]byte(s))
	hash := hex.EncodeToString(h.Sum(nil))
	result := []string{}

	for i := 0; i < _level; i++ {
		result = append(result, hash[i*2:i*2+2])
	}

	return strings.Join(result, "/")
}

func HashPrefix(input string, level ...int) string {
	return filepath.Join(HashSplit(input, level...), input)
}
