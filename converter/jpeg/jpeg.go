package jpeg

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"log/slog"
	"os/exec"
	"path"
	"strings"
)

// Options represent JPEG encoding options.
type Options struct {
	Quality int `json:"quality"`
}

// DecodeJPEG decodes a JPEG file and return an image.
func DecodeJPEG(r io.Reader) (image.Image, string, error) {
	i, realFormat, err := image.Decode(r)
	slog.Info("realFormat", "realFormat", realFormat)
	if err != nil {
		return nil, "", err
	}
	if realFormat == "jpeg" {
		realFormat = "jpg"
	}
	return i, realFormat, nil
}

// EncodeJPEG encodes an image into JPEG and returns a buffer.
func EncodeJPEG(i image.Image, o *Options) (buf bytes.Buffer, err error) {
	err = jpeg.Encode(&buf, i, &jpeg.Options{Quality: o.Quality})
	return buf, err
}

func Encode(inputFile, outDir string) (string, error) {
	slog.Info("Encode JPEG", "inputFile", inputFile)
	if !isJpeg(inputFile) {
		newInputFile := strings.Replace(inputFile, path.Ext(inputFile), ".jpg", 1)
		convertCmd := exec.Command(
			"convert", "-strip",
			inputFile, newInputFile,
		)
		slog.Info("cmd", "cmd", convertCmd.String())
		err := convertCmd.Run()
		if err != nil {
			slog.Error("convert error", "err", err)
			return "", err
		}
		inputFile = newInputFile
	}

	cmd := exec.Command(
		"jpegoptim",
		"--strip-all",
		"-o", "-m", "80",
		inputFile, "-d", outDir,
	)
	slog.Info("cmd", "cmd", cmd.String())
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	outputFile := path.Join(outDir, path.Base(inputFile))
	return outputFile, nil
}

func isJpeg(inputFile string) bool {
	return path.Ext(inputFile) == ".jpg" || path.Ext(inputFile) == ".jpeg"
}
