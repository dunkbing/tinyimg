package image

import (
	"log/slog"
	"optipic/converter/config"
	"optipic/converter/stat"
	"runtime/debug"
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
func (fm *FileManager) Convert() (fileResults []FileResult, zippedUrl string, errs []error) {
	file := fm.File
	fileResults, zippedUrl, errs = file.Write(fm.config)

	for _, f := range fileResults {
		fm.stats.IncreaseByteCount(f.SavedBytes)
		fm.stats.IncreaseTimeCount(f.Time)
		fm.stats.IncreaseImageCount(1)
	}
	fm.Clear()

	return fileResults, zippedUrl, errs
}
