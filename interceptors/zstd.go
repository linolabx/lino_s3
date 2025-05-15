package interceptors

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/klauspost/compress/zstd"
	"github.com/linolabx/lino_s3"
)

var Zstd = &lino_s3.Interceptor{
	PostGet: func(output *s3.GetObjectOutput) (*s3.GetObjectOutput, error) {
		zr, err := zstd.NewReader(output.Body)
		if err != nil {
			return nil, err
		}

		output.Body = zr.IOReadCloser()
		return output, nil
	},
	PrePut: func(input *s3.PutObjectInput) (*s3.PutObjectInput, error) {
		var buf bytes.Buffer
		zz, err := zstd.NewWriter(&buf)
		if err != nil {
			return nil, err
		}

		data, err := io.ReadAll(input.Body)
		if err != nil {
			return nil, err
		}

		if _, err := zz.Write(data); err != nil {
			return nil, err
		}

		if err := zz.Close(); err != nil {
			return nil, err
		}

		input.Body = bytes.NewReader(buf.Bytes())
		return input, nil
	},
	PreUpload: func(input *s3.PutObjectInput) (*s3.PutObjectInput, error) {
		pr, pw := io.Pipe()
		go func() {
			zw, err := zstd.NewWriter(pw)
			if err != nil {
				pw.CloseWithError(err)
				return
			}

			if _, err := io.Copy(zw, input.Body); err != nil {
				pw.CloseWithError(err)
				return
			}

			if err := zw.Close(); err != nil {
				pw.CloseWithError(err)
				return
			}

			pw.Close()
		}()

		input.Body = pr
		return input, nil
	},
}
