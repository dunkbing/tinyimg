package webp

import (
	"image"
	"io"
	"log/slog"
	"os/exec"
	"path"
	"strings"

	"github.com/chai2010/webp"
)

// Options represent WebP encoding options.
type Options struct {
	Lossless bool `json:"lossless"`
	Quality  int  `json:"quality"`
}

// DecodeWebp a webp file and return an image.
func DecodeWebp(r io.Reader) (image.Image, string, error) {
	realFormat := "webp"
	i, err := webp.Decode(r)
	if err != nil {
		i, realFormat, err = image.Decode(r)
		if err != nil {
			return nil, "", err
		}
	}
	return i, realFormat, nil
}

// Encode encodes an image into webp and returns a buffer.
func Encode(inputFile, outDir string) (string, error) {
	slog.Info("Encode WebP", "inputFile", inputFile)
	filename := path.Base(inputFile)
	outputFile := path.Join(outDir, filename)
	outputFile = strings.Replace(outputFile, path.Ext(outputFile), ".webp", 1)

	cmd := exec.Command(
		"cwebp", "-q", "80",
		inputFile, "-o", outputFile,
	)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return outputFile, nil
}
