package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/linolabx/lino_s3/internal"
)

func Hash(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Hash("mykey") => "9adbe0b3033881f88ebd825bcf763b43"
// HashTemplate("mykey", "{.}") => "mykey"
// HashTemplate("mykey", "{hash}") => "9adbe0b3033881f88ebd825bcf763b43"
// HashTemplate("mykey", "{hash.l3}/{.}") => "9a/db/e0/mykey"
// HashTemplate("mykey", "{hash.l4}/{.}") => "9a/db/e0/b3/mykey"
// HashTemplate("mykey", "{hash.l3}/{hash}") => "9a/db/e0/9adbe0b3033881f88ebd825bcf763b43"
func HashTemplate(input string, template string) string {
	hash := Hash(input)
	prefixSlice := []string{}
	for i := range 16 {
		prefixSlice = append(prefixSlice, hash[i*2:i*2+2])
	}
	prefix := strings.Join(prefixSlice, "/")

	result := template
	result = strings.ReplaceAll(result, "{.}", input)
	result = strings.ReplaceAll(result, "{hash}", hash)
	for i := range 16 {
		result = strings.ReplaceAll(result, fmt.Sprintf("{hash.l%d}", i+1), prefix[:i*3+3-1])
	}
	return result
}

// Deprecated: Use HashTemplate instead
func HashSplit(s string, level ...int) string {
	return HashTemplate(s, fmt.Sprintf("{hash.l%d}", internal.OptionalParam(3, level...)))
}

// Deprecated: Use HashTemplate instead
func HashPrefix(s string, level ...int) string {
	return HashTemplate(s, fmt.Sprintf("{hash.l%d}/{.}", internal.OptionalParam(3, level...)))
}
