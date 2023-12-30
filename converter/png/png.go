package png

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"log/slog"

	"github.com/foobaz/lossypng/lossypng"
)

const qMax = 20

// Options represent PNG encoding options.
type Options struct {
	Quality int `json:"quality"`
}

// DecodePNG decodes a PNG file and return an image.
func DecodePNG(r io.Reader) (image.Image, error) {
	i, realFormat, err := image.Decode(r)
	slog.Info("realFormat", "realFormat", realFormat)
	if err != nil {
		return nil, err
	}
	return i, nil
}

// EncodePNG encodes an image into PNG and returns a buffer.
func EncodePNG(i image.Image, o *Options) (buf bytes.Buffer, err error) {
	c := lossypng.Compress(i, 2, qualityFactor(o.Quality))
	err = png.Encode(&buf, c)
	return buf, err
}

// qualityFactor normalizes the PNG quality factor from a max of 20, where 0 is
// no conversion.
func qualityFactor(q int) int {
	f := q / 100
	return qMax - (f * qMax)
}
