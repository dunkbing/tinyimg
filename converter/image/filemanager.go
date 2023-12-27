package image

import (
	"fmt"
	"log/slog"
	"optipic/converter/config"
	"optipic/converter/stat"
	"runtime/debug"
	"strings"
)

// FileManager handles collections of Files for conversion.
type FileManager struct {
	File *File

	Logger *slog.Logger

	config *config.Config
	stats  *stat.Stat
}

// NewFileManager creates a new FileManager.
func NewFileManager() *FileManager {
	logger := slog.Default()
	logger.Info("FileManager initialized...")
	return &FileManager{
		config: config.NewConfig(),
		stats:  stat.NewStat(),
		Logger: logger,
	}
}

// HandleFile processes a file from the client.
func (fm *FileManager) HandleFile(file *File) (err error) {
	if err = file.Decode(); err != nil {
		return err
	}
	fm.File = file
	fm.Logger.Info("added file to file manager", "filename", file.Name)

	return nil
}

// Clear removes the files in the FileManager.
func (fm *FileManager) Clear() {
	fm.File = nil
	debug.FreeOSMemory()
}

// Convert runs the conversion on all files in the FileManager.
func (fm *FileManager) Convert() (stat map[string]any, errs []error) {
	c := 0
	var b int64 // bytes
	file := fm.File
	if !file.IsConverted {
		err := file.Write(fm.config)
		if err != nil {
			fm.Logger.Error("failed to convert file", "fileID", file.ID, "error", err)
			errs = append(errs, fmt.Errorf("failed to convert file: %s", file.Name))
		} else {
			fm.Logger.Info(fmt.Sprintf("converted file: %s", file.Name))
			s, err := file.GetConvertedSize()
			if err != nil {
				fm.Logger.Error("failed to read converted file size", "error", err)
			}
			fm.Logger.Info(
				"converted file",
				"fileID", file.ID,
				"path", strings.Replace(file.ConvertedFile, "\\", "/", -1),
				"size", s,
			)
			c++
			s, err = file.GetSavings()
			if err != nil {
				fm.Logger.Error("failed to get file conversion savings", "error", err)
			}
			b += s
		}
	}

	fm.stats.IncreaseImageCount(c)
	fm.stats.IncreaseByteCount(b)
	fm.stats.IncreaseTimeCount(file.ConvertTime)
	fm.Clear()

	return map[string]any{
		"count":      c,
		"savedBytes": b,
		"time":       file.ConvertTime,
		"imageUrl":   file.S3Url,
	}, errs
}
