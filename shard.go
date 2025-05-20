package lino_s3

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

func Hash(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func ShardT(input string, template string, values ...any) string {
	hash := Hash(input)
	prefixSlice := []string{}
	for i := range 16 {
		prefixSlice = append(prefixSlice, hash[i*2:i*2+2])
	}
	prefix := strings.Join(prefixSlice, "/")

	result := template
	result = strings.ReplaceAll(result, "{.}", input)
	result = strings.ReplaceAll(result, "{hash}", hash)
	result = strings.ReplaceAll(result, "{shard}", "{shard.l3}")
	for i := range 16 {
		result = strings.ReplaceAll(result, fmt.Sprintf("{shard.l%d}", i+1), prefix[:i*3+3-1])
	}
	return fmt.Sprintf(result, values...)
}
