package images

import (
	"encoding/base64"
	"io"
	"os"
)

func ConvertToBase64(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() { err = file.Close() }()

	image, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	defer func() { err = file.Close() }()

	image64 := base64.StdEncoding.EncodeToString(image)

	return image64, err
}
