package image

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dunkbing/tinyimg/converter/cache"
	"github.com/dunkbing/tinyimg/converter/config"
	"github.com/dunkbing/tinyimg/converter/jpeg"
	"github.com/dunkbing/tinyimg/converter/png"
	"github.com/dunkbing/tinyimg/converter/webp"
	"image"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var mimes = map[string]string{
	"image/.jpg": "jpg",
	"image/jpg":  "jpg",
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/webp": "webp",
}

var logger = slog.Default()

// File represents an image file.
type File struct {
	Data          []byte `json:"data"`
	Ext           string `json:"ext"`
	MimeType      string `json:"type"`
	Name          string `json:"name"`
	Size          int64  `json:"size"`
	ConvertedFile string
	InputFileDest string
	Image         image.Image
	Formats       []string
	cache         *cache.Cache[string, CompressResult]
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
		f.Image, f.Ext, err = jpeg.DecodeJPEG(bytes.NewReader(f.Data))
	case "png":
		f.Image, f.Ext, err = png.DecodePNG(bytes.NewReader(f.Data))
	case "webp":
		f.Image, f.Ext, err = webp.DecodeWebp(bytes.NewReader(f.Data))
	default:
		err = errors.New("unsupported file type:" + mime)
	}
	if err != nil {
		return err
	}
	newFileName := strings.Split(f.InputFileDest, ".")[0] + "." + f.Ext
	err = os.Rename(f.InputFileDest, newFileName)
	f.InputFileDest = newFileName
	return err
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

type CompressResult struct {
	SavedBytes int64  `json:"savedBytes"`
	NewSize    int64  `json:"newSize"`
	Time       int64  `json:"time"`
	ImageUrl   string `json:"imageUrl"`
	Format     string `json:"format"`
}

// Write saves a file to disk based on the encoding target.
func (f *File) Write(c *config.Config) ([]CompressResult, []string, []error) {
	var errs []error
	t := time.Now()
	var compressedFiles []string

	formats := f.Formats
	res := make([]CompressResult, len(formats))
	var wg sync.WaitGroup
	wg.Add(len(formats))
	for i, format := range formats {
		go func(format string, index int) {
			defer wg.Done()
			var savedBytes, newSize int64 // bytes
			filename := strings.Split(f.Name, ".")[0]
			filename = filename + "." + format

			compressedFiles = append(compressedFiles, filename)
			if cachedRes, ok := f.cache.Get(filename); ok {
				res[index] = cachedRes
				return
			}

			outputFile, err := encToBuf(f, format)
			if err != nil {
				errs = append(errs, err)
				return
			}

			nt := time.Since(t).Milliseconds()

			f.ConvertedFile = filepath.Clean(outputFile)
			savedBytes, _ = f.GetSavings()
			newSize, _ = f.GetConvertedSize()
			imageUrl := fmt.Sprintf("%s/image?f=%s", config.HostUrl, filename)

			res[index] = CompressResult{
				SavedBytes: savedBytes,
				NewSize:    newSize,
				Time:       nt,
				ImageUrl:   imageUrl,
				Format:     format,
			}
			f.cache.Set(filename, res[index])
		}(format, i)
	}
	wg.Wait()

	return res, compressedFiles, errs
}

// encToBuf encodes an image to a buffer using the configured target.
func encToBuf(f *File, target string) (outputFile string, err error) {
	c := config.GetConfig()
	switch target {
	case "jpg", "jpeg":
		outputFile, err = jpeg.Encode(f.InputFileDest, c.App.OutDir)
	case "png":
		outputFile, err = png.Encode(f.InputFileDest, c.App.OutDir)
	case "webp":
		outputFile, err = webp.Encode(f.InputFileDest, c.App.OutDir)
	}
	if err != nil {
		return "", err
	}
	return outputFile, nil
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
		err := func(filename string) error {
			f, err := os.Open(path.Join(c.App.OutDir, file))
			if err != nil {
				return err
			}
			defer f.Close()

			info, err := f.Stat()
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			header.Name = filepath.Base(file)
			header.Method = zip.Deflate

			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			_, err = io.Copy(writer, f)
			if err != nil {
				return err
			}
			return nil
		}(file)

		if err != nil {
			continue
		}
	}

	nt := time.Since(t).Milliseconds()
	logger.Info("zipped files", "name", name, "time", nt)

	return zipFile.Name(), nil
}
