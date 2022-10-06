package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
)

func ByteToStdBase64String(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func ImageToPngBytes(image image.Image) ([]byte, error) {
	var pngBuf bytes.Buffer
	err := png.Encode(&pngBuf, image)
	if err != nil {
		return nil, err
	}
	return pngBuf.Bytes(), nil
}

// ImageToPngBase64String
// 1. image to png bytes
// 2. png bytes to base64 string
func ImageToPngBase64String(image image.Image) (string, error) {
	pngBytes, err := ImageToPngBytes(image)
	if err != nil {
		return "", err
	}
	return ByteToStdBase64String(pngBytes), nil
}
