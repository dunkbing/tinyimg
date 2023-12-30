package config

import (
	"fmt"
	"optipic/converter/jpeg"
	"optipic/converter/png"
	"optipic/converter/webp"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

const filename = "conf.json"

// App represents application persistent configuration values.
type App struct {
	OutDir  string        `json:"outDir"`
	Target  string        `json:"target"`
	Sizes   []*size       `json:"sizes"`
	JpegOpt *jpeg.Options `json:"jpegOpt"`
	PngOpt  *png.Options  `json:"pngOpt"`
	WebpOpt *webp.Options `json:"webpOpt"`
}

// Config represents the application settings.
type Config struct {
	App *App
}

// NewConfig returns a new instance of Config.
func NewConfig() *Config {
	c := &Config{}
	c.App, _ = defaults()

	return c
}

// GetAppConfig returns the application configuration.
func (c *Config) GetAppConfig() map[string]interface{} {
	return map[string]interface{}{
		"outDir":  c.App.OutDir,
		"target":  c.App.Target,
		"sizes":   c.App.Sizes,
		"jpegOpt": c.App.JpegOpt,
		"pngOpt":  c.App.PngOpt,
		"webpOpt": c.App.WebpOpt,
	}
}

// RestoreDefaults sets the app configuration to defaults.
func (c *Config) RestoreDefaults() (err error) {
	var a *App
	a, err = defaults()
	if err != nil {
		return err
	}
	c.App = a
	return nil
}

// defaults returns the application configuration defaults.
func defaults() (*App, error) {
	a := &App{
		Target:  "webp",
		JpegOpt: &jpeg.Options{Quality: 80},
		PngOpt:  &png.Options{Quality: 80},
		WebpOpt: &webp.Options{Lossless: false, Quality: 80},
	}
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get user directory: %v", err)
		return nil, err
	}

	od := path.Join(wd, "output")
	cp := filepath.Clean(od)

	if _, err = os.Stat(od); os.IsNotExist(err) {
		if err = os.Mkdir(od, 0777); err != nil {
			od = "./"
			fmt.Printf("failed to create default output directory: %v", err)
			return nil, err
		}
	}
	a.OutDir = cp
	return a, nil
}

// rect represents an image width and height size.
type rect struct {
	Height int `json:"height,omitempty"`
	Width  int `json:"width,omitempty"`
}

// String returns a string representation of the rect.
// For example, "1280x720"
func (r *rect) String() string {
	w := strconv.Itoa(r.Width)
	h := strconv.Itoa(r.Height)
	return fmt.Sprintf("%sx%s", w, h)
}

// size represents an image resizing. Strategy represents an image resizing
// strategy, such as cropping.
type size struct {
	rect
	Strategy int `json:"strategy"`
}
