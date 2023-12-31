package image

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"io"
	"log/slog"
	"optipic/converter/config"
	"optipic/converter/jpeg"
	"optipic/converter/png"
	"optipic/converter/webp"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
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
	Image         image.Image
	Formats       []string
}

// Decode decodes the file's data based on its mime type.
func (f *File) Decode() error {
	mime, err := GetFileType(f.MimeType)
	logger.Info("mime", "mime", mime)
	if err != nil {
		return err
	}

	switch mime {
	case "jpg", "jpeg":
		f.Image, err = jpeg.DecodeJPEG(bytes.NewReader(f.Data))
	case "png":
		f.Image, err = png.DecodePNG(bytes.NewReader(f.Data))
	case "webp":
		f.Image, err = webp.DecodeWebp(bytes.NewReader(f.Data))
	default:
		err = errors.New("unsupported file type:" + mime)
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

type FileResult struct {
	SavedBytes int64  `json:"savedBytes"`
	NewSize    int64  `json:"newSize"`
	Time       int64  `json:"time"`
	ImageUrl   string `json:"imageUrl"`
	Format     string `json:"format"`
}

// Write saves a file to disk based on the encoding target.
func (f *File) Write(c *config.Config) ([]FileResult, []string, []error) {
	// TODO resizing should probably be in its own method
	var errs []error
	t := time.Now()
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
					errs = append(errs, err)
					continue
				}
				croppedImg := f.Image.(SubImager).SubImage(crop)
				i = imaging.Resize(croppedImg, r.Width, r.Height, imaging.Lanczos)
				s = fmt.Sprintf("%dx%d", i.Bounds().Max.X, i.Bounds().Max.Y)
			}
			buf, err := encToBuf(i, f.Ext)
			dest := path.Join(c.App.OutDir, f.Name+"--"+s+"."+f.Ext)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if err = os.WriteFile(dest, buf.Bytes(), 0666); err != nil {
				errs = append(errs, err)
				continue
			}
		}
	}
	compressedFiles := []string{}
	s3Client, err := NewS3Client()
	if err != nil {
		logger.Error("failed to create s3 client", "error", err)
	}
	formats := f.Formats
	res := make([]FileResult, len(formats))
	var wg sync.WaitGroup
	wg.Add(len(formats))
	for i, format := range formats {
		fmt.Println("format", format)
		go func(format string, index int) {
			defer wg.Done()
			var savedBytes, newSize int64 // bytes
			buf, err := encToBuf(f.Image, format)
			if err != nil {
				errs = append(errs, err)
				return
			}
			filename := strings.Split(f.Name, ".")[0]
			filename = filename + "." + format
			dest := path.Join(c.App.OutDir, filename)
			compressedFiles = append(compressedFiles, filename)
			if err = os.WriteFile(dest, buf.Bytes(), 0666); err != nil {
				errs = append(errs, err)
				return
			}
			nt := time.Since(t).Milliseconds()

			err = s3Client.UploadFile(filename, dest)
			if err != nil {
				logger.Error("failed to upload file to s3", "error", err)
			}
			f.ConvertedFile = filepath.Clean(dest)
			savedBytes, _ = f.GetSavings()
			newSize, _ = f.GetConvertedSize()
			imageUrl, _ := s3Client.GetFileUrl(filename)

			res[index] = FileResult{
				SavedBytes: savedBytes,
				NewSize:    newSize,
				Time:       nt,
				ImageUrl:   imageUrl,
				Format:     format,
			}
		}(format, i)
	}
	wg.Wait()

	return res, compressedFiles, errs
}

// encToBuf encodes an image to a buffer using the configured target.
func encToBuf(i image.Image, target string) (*bytes.Buffer, error) {
	var b bytes.Buffer
	var err error
	switch target {
	case "jpg", "jpeg":
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

// GetFileType returns the file's type based on the given mime type.
func GetFileType(t string) (string, error) {
	m, prs := mimes[t]
	if !prs {
		_ = errors.New("unsupported file type:" + t)
	}
	return m, nil
}

func generateUniqueZipFilename(files []string) string {
	// Concatenate all file names
	var concatenatedNames string
	for _, file := range files {
		concatenatedNames += filepath.Base(file)
	}

	// Calculate MD5 hash
	hash := md5.New()
	hash.Write([]byte(concatenatedNames))
	hashInBytes := hash.Sum(nil)

	// Convert the hash to a hexadecimal string
	hashString := hex.EncodeToString(hashInBytes)

	// Use the hash as the zip filename with a ".zip" extension
	return hashString + ".zip"
}

func zipFiles(files []string, c *config.Config) (string, error) {
	logger.Info("zipping files", "files", files)
	t := time.Now()
	baseFiles := []string{}
	for _, file := range files {
		baseFiles = append(baseFiles, filepath.Base(file))
	}
	name := generateUniqueZipFilename(baseFiles)
	zipFile, err := os.Create(name)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		f, err := os.Open(path.Join(c.App.OutDir, file))
		if err != nil {
			continue
		}
		defer f.Close()

		info, err := f.Stat()
		if err != nil {
			continue
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			continue
		}

		header.Name = filepath.Base(file)
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			continue
		}

		_, err = io.Copy(writer, f)
		if err != nil {
			return "", err
		}
	}

	nt := time.Since(t).Milliseconds()
	logger.Info("zipped files", "name", name, "time", nt)

	return zipFile.Name(), nil
}

// SubImager handles creating a subimage from an image rect.
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}
