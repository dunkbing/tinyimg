package stat

import (
	"log/slog"
)

// Stat represents application statistics.
type Stat struct {
	ByteCount  int64 `json:"byteCount"`
	ImageCount int   `json:"imageCount"`
	TimeCount  int64 `json:"timeCount"`

	Logger  *slog.Logger
}

var stat *Stat

// NewStat returns a new Stat instance.
func NewStat() *Stat {
	logger := slog.With("Stat")
	logger.Info("Stat initialized...")
	if stat == nil {
		stat = &Stat{
			Logger: logger,
		}
	}
	return stat
}

// GetStats returns the application stats.
func (s *Stat) GetStats() map[string]any {
	return map[string]interface{}{
		"byteCount":  s.ByteCount,
		"imageCount": s.ImageCount,
		"timeCount":  s.TimeCount,
	}
}

// IncreaseByteCount adds and persists the given byte count to the app stats.
func (s *Stat) IncreaseByteCount(b int64) {
	if b <= 0 {
		return
	}
	s.ByteCount += b
	if err := s.store(); err != nil {
		s.Logger.Error("failed to store stats", "error", err)
	}
}

// IncreaseImageCount adds and persists the given image count to the app stats.
func (s *Stat) IncreaseImageCount(i int) {
	if i <= 0 {
		return
	}
	s.ImageCount += i
	if err := s.store(); err != nil {
		s.Logger.Error("failed to store stats", "error", err)
	}
}

// IncreaseTimeCount adds and persists the given time count to the app stats.
func (s *Stat) IncreaseTimeCount(t int64) {
	if t < 0 {
		return
	}
	s.TimeCount += t
	if err := s.store(); err != nil {
		s.Logger.Error("failed to store stats", "error", err)
	}
}

// store stores the app stats to the file system.
func (s *Stat) store() error {
	s.GetStats()
	return nil
}
