package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
)

func UploadFile(file multipart.File) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, fmt.Errorf("could not upload file. %v", err)
	}

	return buf.Bytes(), nil
}
