package utils

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GenerateHash(str string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	hashStr := fmt.Sprintf("%x", hashBytes)

	return hashStr, nil
}

func IsValidUrl(url_ string) bool {
	_, err := url.ParseRequestURI(url_)
	if err != nil {
		return false
	}
	return true
}

func ZipFolder(src, dest string) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	zipWriter := zip.NewWriter(out)
	defer zipWriter.Close()

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate // Use Deflate compression

		header.Name, err = filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			header.Name += "/" // Add trailing slash for directories
		} else {
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func GetFirstFile(dirPath string) (string, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return "", err
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() {
			return filepath.Join(dirPath, file.Name()), nil
		}
	}

	return "", os.ErrNotExist
}

func DownloadVideo(url, outDir string) (string, error) {
	// Create the new directory if it doesn't exist
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	slog.Info("Downloading video", "url", url)
	id := uuid.New().String()
	outDest := filepath.Join(outDir, id)
	output := outDest + "/%(title)s.%(ext)s"
	args := []string{
		"-o", output, "--quiet", url,
	}
	if strings.Contains(url, "tiktok") {
		args = append([]string{"-f", "0"}, args...)
	}
	cmd := exec.Command("yt-dlp", args...)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error executing yt-dlp: %w", err)
	}

	// get file name
	downloadedFilePath, err := GetFirstFile(outDest)
	if err != nil {
		slog.Error("Error getting downloaded file", "error", err)
		return "", err
	}
	filename := filepath.Base(downloadedFilePath)
	ext := filepath.Ext(filename)
	newFilename := uuid.New().String()
	newFilename = fmt.Sprintf("%s%s", newFilename, ext)
	filepath_ := filepath.Join(outDir, newFilename)
	err = os.Rename(downloadedFilePath, filepath_)
	if err != nil {
		return "", err
	}

	return newFilename, nil
}

func DownloadPlaylist(url, outDir string) (string, error) {
	// Create the new directory if it doesn't exist
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	slog.Info("Downloading video", "url", url)
	id := uuid.New().String()
	outDest := filepath.Join(outDir, id)
	output := outDest + "/%(playlist_index)s - %(title)s.%(ext)s"
	cmd := exec.Command("yt-dlp", "-o", output, "--quiet", url)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error executing yt-dlp: %w", err)
	}

	// get file name
	zipFile := fmt.Sprintf("%s.zip", id)
	zipDest := fmt.Sprintf("%s.zip", outDest)
	err = ZipFolder(outDest, zipDest)
	if err != nil {
		return "", err
	}

	return zipFile, nil
}
