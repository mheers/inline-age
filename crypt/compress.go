package crypt

import (
	"bytes"
	"encoding/base64"
	"io"

	"github.com/andybalholm/brotli"
)

// compressBrotli compresses a byte slice using brotli compression.
func compressBrotli(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := brotli.NewWriterLevel(&b, brotli.BestCompression)
	w.Write(data)
	w.Close()
	return b.Bytes(), nil
}

// decompressBrotli decompresses a byte slice using brotli compression.
func decompressBrotli(data []byte) ([]byte, error) {
	r := brotli.NewReader(bytes.NewReader(data))
	result, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// capsule compresses a byte slice using brotli compression and returns a base64 encoded string.
func capsule(in []byte) (string, error) {
	compressed, err := compressBrotli(in)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(compressed), nil
}

// uncapsule decodes a base64 encoded string and decompresses it using brotli compression.
func uncapsule(in string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return nil, err
	}

	return decompressBrotli(decoded)
}
