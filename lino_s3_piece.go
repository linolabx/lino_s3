package lino_s3

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type LinoS3Piece struct {
	sess   *session.Session
	bucket string
	key    string

	start int64
	end   int64
}

func (s *LinoS3Piece) Get() (*s3.GetObjectOutput, error) {
	Range := fmt.Sprintf("bytes=%d-%d", s.start, s.end-1)

	out, err := s3.New(s.sess).GetObject(&s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &s.key,
		Range:  &Range,
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *LinoS3Piece) ReadTo(writer io.WriteCloser) error {
	out, err := s.Get()
	if err != nil {
		return err
	}

	defer out.Body.Close()
	defer writer.Close()

	_, err = io.Copy(writer, out.Body)
	return err
}

func (s *LinoS3Piece) ReadBuffer() ([]byte, error) {
	out, err := s.Get()
	if err != nil {
		return nil, err
	}

	defer out.Body.Close()
	return io.ReadAll(out.Body)
}

func (s *LinoS3Piece) ReadString() (string, error) {
	out, err := s.Get()
	if err != nil {
		return "", err
	}

	defer out.Body.Close()
	data, err := io.ReadAll(out.Body)
	return string(data), err
}

func (s *LinoS3Piece) ReadJSON(v interface{}) error {
	data, err := s.ReadBuffer()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func (s *LinoS3Piece) ReadCBOR(v interface{}) error {
	data, err := s.ReadBuffer()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func (s *LinoS3Piece) Key() string {
	return fmt.Sprintf("%s?range=%d-%d", s.key, s.start, s.end)
}
