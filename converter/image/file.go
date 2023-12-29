package image

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"log/slog"
	"optipic/converter/config"
	"optipic/converter/jpeg"
	"optipic/converter/png"
	"optipic/converter/webp"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
)

const (
	fill = iota
	fit
	smart
)

var mimes = map[string]string{
	"image/.jpg": "jpg",
	"image/jpg":  "jpg",
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/webp": "webp",
}

var logger *slog.Logger = slog.Default()

// File represents an image file.
type File struct {
	Data          []byte `json:"data"`
	Ext           string `json:"ext"`
	MimeType      string `json:"type"`
	Name          string `json:"name"`
	Size          int64  `json:"size"`
	ConvertedFile string
	IsConverted   bool
	ConvertTime   int64
	S3Url         string
	Image         image.Image
}

// Decode decodes the file's data based on its mime type.
func (f *File) Decode() error {
	mime, err := getFileType(f.MimeType)
	if err != nil {
		return err
	}

	switch mime {
	case "jpg":
		f.Image, err = jpeg.DecodeJPEG(bytes.NewReader(f.Data))
	case "png":
		f.Image, err = png.DecodePNG(bytes.NewReader(f.Data))
	case "webp":
		f.Image, err = webp.DecodeWebp(bytes.NewReader(f.Data))
	}
	if err != nil {
		return err
	}
	return nil
}

// GetConvertedSize returns the size of the converted file.
func (f *File) GetConvertedSize() (int64, error) {
	if f.ConvertedFile == "" {
		return 0, errors.New("file has no converted file")
	}
	s, err := os.Stat(f.ConvertedFile)
	if err != nil {
		return 0, err
	}
	return s.Size(), nil
}

// GetSavings returns the delta between original and converted file size.
func (f *File) GetSavings() (int64, error) {
	c, err := f.GetConvertedSize()
	if err != nil {
		return 0, err
	}
	return f.Size - c, nil
}

// Write saves a file to disk based on the encoding target.
func (f *File) Write(c *config.Config) error {
	// TODO resizing should probably be in its own method
	t := time.Now().UnixNano()
	if c.App.Sizes != nil {
		for _, r := range c.App.Sizes {
			if r.Height <= 0 || r.Width <= 0 {
				logger.Warn("invalid image size", "size", r.String())
				continue
			}
			var i image.Image
			var s string
			switch r.Strategy {
			case fill:
				i = imaging.Fill(f.Image, r.Width, r.Height, imaging.Center, imaging.Lanczos)
				s = r.String()
			case fit:
				i = imaging.Fit(f.Image, r.Width, r.Height, imaging.Lanczos)
				s = fmt.Sprintf("%dx%d", i.Bounds().Max.X, i.Bounds().Max.Y)
			case smart:
				analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
				crop, err := analyzer.FindBestCrop(f.Image, r.Width, r.Height)
				if err != nil {
					return err
				}
				croppedImg := f.Image.(SubImager).SubImage(crop)
				i = imaging.Resize(croppedImg, r.Width, r.Height, imaging.Lanczos)
				s = fmt.Sprintf("%dx%d", i.Bounds().Max.X, i.Bounds().Max.Y)
			}
			buf, err := encToBuf(i, f.Ext)
			dest := path.Join(c.App.OutDir, c.App.Prefix+f.Name+"--"+s+c.App.Suffix+"."+f.Ext)
			if err != nil {
				return err
			}
			if err = os.WriteFile(dest, buf.Bytes(), 0666); err != nil {
				return err
			}
		}
	}
	buf, err := encToBuf(f.Image, c.App.Target)

	filename := strings.Split(f.Name, ".")[0]
	dest := path.Join(c.App.OutDir, c.App.Prefix+filename+c.App.Suffix+"."+c.App.Target)
	if err != nil {
		return err
	}
	logger.Info("writing file", "path", dest)
	if err = os.WriteFile(dest, buf.Bytes(), 0666); err != nil {
		return err
	}
	nt := (time.Now().UnixNano() - t) / 1000000 // milliseconds
	f.ConvertTime = nt
	s3Client, err := NewS3Client()
	if err != nil {
		logger.Error("failed to create s3 client", "error", err)
	}
	err = s3Client.UploadFile(f.Name, dest)
	if err != nil {
		logger.Error("failed to upload file to s3", "error", err)
	}
	f.ConvertedFile = filepath.Clean(dest)
	f.IsConverted = true
	f.S3Url = s3Client.GetFileUrl(f.Name)
	return nil
}

// encToBuf encodes an image to a buffer using the configured target.
func encToBuf(i image.Image, target string) (*bytes.Buffer, error) {
	var b bytes.Buffer
	var err error
	switch target {
	case "jpg":
		b, err = jpeg.EncodeJPEG(i, &jpeg.Options{Quality: 80})
	case "png":
		b, err = png.EncodePNG(i, &png.Options{Quality: 80})
	case "webp":
		b, err = webp.EncodeWebp(i, &webp.Options{Lossless: false, Quality: 80})
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// getFileType returns the file's type based on the given mime type.
func getFileType(t string) (string, error) {
	m, prs := mimes[t]
	if !prs {
		_ = errors.New("unsupported file type:" + t)
	}
	return m, nil
}

// SubImager handles creating a subimage from an image rect.
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}
