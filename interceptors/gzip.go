package interceptors

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/linolabx/lino_s3"
)

var Gzip = &lino_s3.Interceptor{
	PostGet: func(output *s3.GetObjectOutput) (*s3.GetObjectOutput, error) {
		gr, err := gzip.NewReader(output.Body)
		if err != nil {
			return nil, err
		}

		output.Body = gr
		return output, nil
	},
	PrePut: func(input *s3.PutObjectInput) (*s3.PutObjectInput, error) {
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)

		data, err := io.ReadAll(input.Body)
		if err != nil {
			return nil, err
		}

		if _, err := gz.Write(data); err != nil {
			return nil, err
		}

		if err := gz.Close(); err != nil {
			return nil, err
		}

		input.Body = bytes.NewReader(buf.Bytes())
		return input, nil
	},
	PreUpload: func(input *s3.PutObjectInput) (*s3.PutObjectInput, error) {
		pr, pw := io.Pipe()
		go func() {
			gw := gzip.NewWriter(pw)

			if _, err := io.Copy(gw, input.Body); err != nil {
				pw.CloseWithError(err)
				return
			}

			if err := gw.Close(); err != nil {
				pw.CloseWithError(err)
				return
			}

			pw.Close()
		}()

		input.Body = pr
		return input, nil
	},
}
