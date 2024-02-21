package lino_s3

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go/aws/session"
)

type LinoS3Bucket struct {
	sess   *session.Session
	bucket string

	interceptors []*Interceptor
}

func (s *LinoS3Bucket) UseInterceptors(interceptors ...*Interceptor) *LinoS3Bucket {
	s.interceptors = interceptors
	return s
}

func (s *LinoS3Bucket) SubPath(path string) *LinoS3Path {
	return &LinoS3Path{
		sess:         s.sess,
		bucket:       s.bucket,
		path:         path,
		interceptors: s.interceptors,
	}
}

func (s *LinoS3Bucket) Object(key string) *LinoS3Object {
	return &LinoS3Object{
		sess:         s.sess,
		bucket:       s.bucket,
		key:          key,
		interceptors: s.interceptors,
	}
}

func (s *LinoS3Bucket) Piece(key string, start int64, end int64) *LinoS3Piece {
	return &LinoS3Piece{
		sess:   s.sess,
		bucket: s.bucket,
		key:    key,
		start:  start,
		end:    end,
	}
}

var pieceKeyRegex = regexp.MustCompile(`^(.*)\?range=(\d+)-(\d+)$`)

func mustParseInt(str string) int64 {
	result, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return result
}

func (s *LinoS3Bucket) PieceByKey(pieceKey string) (*LinoS3Piece, error) {
	result := pieceKeyRegex.FindStringSubmatch(pieceKey)
	if result == nil {
		return nil, fmt.Errorf("invalid piece key: %s", pieceKey)
	}

	return s.Piece(result[1], mustParseInt(result[2]), mustParseInt(result[3])), nil
}
