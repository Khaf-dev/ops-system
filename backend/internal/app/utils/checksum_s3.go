package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"
)

func CalculateCheckSum(f *multipart.FileHeader) (string, []byte, error) {
	src, err := f.Open()
	if err != nil {
		return "", nil, err
	}
	defer src.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, src)
	if err != nil {
		return "", nil, err
	}

	hash := md5.Sum(buf.Bytes())
	return hex.EncodeToString(hash[:]), buf.Bytes(), nil
}

func UploadToS3(data []byte, filename, contentType string) (string, error) {
	// TODO : integrate with AWS SDK S3
	// sementara return null dummy supaya system bisa jalan
	return "https://fake-s3.local" + filename, nil
}
