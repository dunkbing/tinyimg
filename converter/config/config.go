package config

import (
	"fmt"
	"github.com/dunkbing/tinyimg/converter/jpeg"
	"github.com/dunkbing/tinyimg/converter/png"
	"github.com/dunkbing/tinyimg/converter/webp"
	"os"
	"path"
	"path/filepath"
)

var AllowedOrigin = os.Getenv("ALLOWED_ORIGIN")
var HostUrl = os.Getenv("HOST_URL")

// App represents application persistent configuration values.
type App struct {
	InDir   string        `json:"inDir"`
	OutDir  string        `json:"outDir"`
	Target  string        `json:"target"`
	JpegOpt *jpeg.Options `json:"jpegOpt"`
	PngOpt  *png.Options  `json:"pngOpt"`
	WebpOpt *webp.Options `json:"webpOpt"`
}

// Config represents the application settings.
type Config struct {
	App *App
}

var config_ *Config

// GetConfig returns a new instance of Config.
func GetConfig() *Config {
	if config_ == nil {
		c := &Config{}
		c.App, _ = defaults()
		config_ = c
	}

	return config_
}

// GetAppConfig returns the application configuration.
func (c *Config) GetAppConfig() map[string]interface{} {
	return map[string]interface{}{
		"inDir":   c.App.InDir,
		"outDir":  c.App.OutDir,
		"target":  c.App.Target,
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
	wd, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("failed to get user directory: %v", err)
		return nil, err
	}

	od := path.Join(wd, "tinyimg", "output")
	cp := filepath.Clean(od)

	if _, err = os.Stat(od); os.IsNotExist(err) {
		if err = os.MkdirAll(od, 0777); err != nil {
			od = "./"
			fmt.Printf("failed to create default output directory: %v", err)
			return nil, err
		}
	}
	a.OutDir = cp

	id := path.Join(wd, "tinyimg", "input")
	cip := filepath.Clean(id)

	if _, err = os.Stat(id); os.IsNotExist(err) {
		if err = os.MkdirAll(id, 0777); err != nil {
			id = "./"
			fmt.Printf("failed to create default input directory: %v", err)
			return nil, err
		}
	}
	a.InDir = cip

	return a, nil
}
