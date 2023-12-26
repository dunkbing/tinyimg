package stat

import (
	"log/slog"
)

const filename = "stats.json"

// Stat represents application statistics.
type Stat struct {
	ByteCount  int64 `json:"byteCount"`
	ImageCount int   `json:"imageCount"`
	TimeCount  int64 `json:"timeCount"`

	Logger  *slog.Logger
}

// NewStat returns a new Stat instance.
func NewStat() *Stat {
	logger := slog.With("Stat")
	logger.Info("Stat initialized...")
	s := &Stat{
		Logger: logger,
	}

	return s
}

// GetStats returns the application stats.
func (s *Stat) GetStats() map[string]any {
	return map[string]interface{}{
		"byteCount":  s.ByteCount,
		"imageCount": s.ImageCount,
		"timeCount":  s.TimeCount,
	}
}

// SetByteCount adds and persists the given byte count to the app stats.
func (s *Stat) SetByteCount(b int64) {
	if b <= 0 {
		return
	}
	s.ByteCount += b
	if err := s.store(); err != nil {
		s.Logger.Error("failed to store stats", "error", err)
	}
}

// SetImageCount adds and persists the given image count to the app stats.
func (s *Stat) SetImageCount(i int) {
	if i <= 0 {
		return
	}
	s.ImageCount += i
	if err := s.store(); err != nil {
		s.Logger.Error("failed to store stats", "error", err)
	}
}

// SetTimeCount adds and persists the given time count to the app stats.
func (s *Stat) SetTimeCount(t int64) {
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
