package png

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"log/slog"
	"os/exec"
	"path"
	"strings"

	"github.com/foobaz/lossypng/lossypng"
)

const qMax = 20

// Options represent PNG encoding options.
type Options struct {
	Quality int `json:"quality"`
}

// DecodePNG decodes a PNG file and return an image.
func DecodePNG(r io.Reader) (image.Image, string, error) {
	i, realFormat, err := image.Decode(r)
	slog.Info("realFormat", "realFormat", realFormat)
	if err != nil {
		return nil, "", err
	}
	return i, realFormat, nil
}

// EncodePNG encodes an image into PNG and returns a buffer.
func EncodePNG(i image.Image, o *Options) (buf bytes.Buffer, err error) {
	c := lossypng.Compress(i, 2, qualityFactor(o.Quality))
	err = png.Encode(&buf, c)
	return buf, err
}

func Encode(inputFile, outDir string) (string, error) {
	slog.Info("Encode PNG", "inputFile", inputFile)
	filename := path.Base(inputFile)
	if !isPng(inputFile) {
		newInputFile := strings.Replace(inputFile, path.Ext(inputFile), ".png", 1)
		convertCmd := exec.Command(
			"vips", "copy",
			inputFile, fmt.Sprintf("%s[strip]", newInputFile),
		)
		err := convertCmd.Run()
		if err != nil {
			slog.Error("convert to png error", "err", err)
			return "", err
		}
		inputFile = newInputFile
	}
	outputFile := path.Join(outDir, filename)
	outputFile = strings.Replace(outputFile, path.Ext(outputFile), ".png", 1)

	cmd := exec.Command(
		"pngquant", "--quality=0-80",
		"--speed=1", inputFile,
		"--output", outputFile,
		"--force", "--strip",
	)
	err := cmd.Run()
	if err != nil {
		slog.Error("pngquant error", "err", err, "command", cmd.String())
		return "", err
	}

	return outputFile, nil
}

// qualityFactor normalizes the PNG quality factor from a max of 20, where 0 is
// no conversion.
func qualityFactor(q int) int {
	f := q / 100
	return qMax - (f * qMax)
}

func isPng(inputFile string) bool {
	return path.Ext(inputFile) == ".png"
}
