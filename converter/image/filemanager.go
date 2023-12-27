package image

import (
	"fmt"
	"log/slog"
	"optipic/converter/config"
	"optipic/converter/stat"
	"runtime/debug"
	"sync"
	"time"
)

// FileManager handles collections of Files for conversion.
type FileManager struct {
	Files []*File

	Logger  *slog.Logger

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
	fm.Files = append(fm.Files, file)
	fm.Logger.Info("added file to file manager", "filename", file.Name)

	return nil
}

// Clear removes the files in the FileManager.
func (fm *FileManager) Clear() {
	fm.Files = nil
	debug.FreeOSMemory()
}

// Convert runs the conversion on all files in the FileManager.
func (fm *FileManager) Convert() (stat map[string]any, errs []error) {
	var wg sync.WaitGroup
	wg.Add(fm.countUnconverted())

	c := 0
	var b int64
	t := time.Now().UnixNano()
	for _, file := range fm.Files {
		file := file
		if !file.IsConverted {
			go func(wg *sync.WaitGroup) {
				err := file.Write(fm.config)
				if err != nil {
					fm.Logger.Error("failed to convert file", "fileID", file.ID, "error", err)
					// fm.Runtime.Events.Emit("notify", map[string]interface{}{
					// 	"msg":  fmt.Sprintf("Failed to convert file: %s, %s", file.Name, err.Error()),
					// 	"type": "warn",
					// })
					errs = append(errs, fmt.Errorf("failed to convert file: %s", file.Name))
				} else {
					fm.Logger.Info(fmt.Sprintf("converted file: %s", file.Name))
					s, err := file.GetConvertedSize()
					if err != nil {
						fm.Logger.Error("failed to read converted file size", "error", err)
					}
					// fm.Runtime.Events.Emit("conversion:complete", map[string]interface{}{
					// 	"id": file.ID,
					// 	// TODO: standardize this path conversion
					// 	"path": strings.Replace(file.ConvertedFile, "\\", "/", -1),
					// 	"size": s,
					// })
					c++
					s, err = file.GetSavings()
					if err != nil {
						fm.Logger.Error("failed to get file conversion savings", "error", err)
					}
					b += s
				}
				wg.Done()
			}(&wg)
		}
	}

	wg.Wait()
	nt := (time.Now().UnixNano() - t) / 1000000
	fm.stats.IncreaseImageCount(c)
	fm.stats.IncreaseByteCount(b)
	fm.stats.IncreaseTimeCount(nt)
	fm.Clear()

	return map[string]any{
		"count":   c,
		"resizes": c * len(fm.config.App.Sizes),
		"savings": b,
		"time":    nt,
	}, errs
}

// countUnconverted returns the number of files in the FileManager that haven't
// been converted.
func (fm *FileManager) countUnconverted() int {
	c := 0
	for _, file := range fm.Files {
		if !file.IsConverted {
			c++
		}
	}
	return c
}
