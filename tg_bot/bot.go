package tg_bot

import (
	"context"
	"fmt"
	"github.com/dunkbing/tinyimg/converter/config"
	"github.com/dunkbing/tinyimg/converter/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
)

type handler struct {
	config *config.Config
}

func New() *handler {
	c := config.GetConfig()
	return &handler{
		config: c,
	}
}

func (h *handler) Start() {
	slog.Info("Starting Telegram Bot")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(h.helpHandler),
	}
	botToken := config.TgBotToken

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}
	_, err = b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{
				Command:     "start",
				Description: "Help",
			},
			{
				Command:     "download",
				Description: "Download video with the url",
			},
		},
	})
	if err != nil {
		return
	}
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypeExact,
		h.helpHandler,
	)
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/download",
		bot.MatchTypeExact,
		h.defaultHandler,
	)

	b.Start(ctx)
}

func (h *handler) helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	help := `
Welcome to the Video Downloader Bot!

This bot allows you to download videos and playlists from various sources.

📥 /video <url>
Download a single video by providing the video URL.

📥 /playlist <url>
Download an entire playlist by providing the playlist URL.
The bot will download all videos in the playlist.

Example usage:
/video https://example.com/video.mp4
/playlist https://example.com/playlist

Note: Video files will be sent as attachments in the chat.

For any issues or feedback, please contact the bot owner.

`
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   help,
	})
	if err != nil {
		slog.Error("Error sending help message", "err", err)
	}
}

func (h *handler) videoHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	url_ := update.Message.Text
	if strings.HasPrefix(url_, "/download ") {
		url_ = strings.Replace(url_, "/download ", "", 1)
	}
	if !isValidUrl(url_) {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Invalid URL",
		})
		return
	}
	filename, err := downloadVideo(url_, h.config.App.OutDir)
	if err != nil {
		slog.Error("Error downloading video", "error", err)
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Something went wrong",
		})
		return
	}
	slog.Info("Downloaded file", "filename", filename)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("%s/video?f=%s", config.HostUrl, filename),
	})

	if err != nil {
		slog.Error("send message", slog.String("error", err.Error()))
		return
	}
}

func isValidUrl(url_ string) bool {
	_, err := url.ParseRequestURI(url_)
	if err != nil {
		return false
	}
	return true
}

func downloadVideo(url, outDir string) (string, error) {
	// Create the new directory if it doesn't exist
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	slog.Info("Downloading video", "url", url)
	cmd := exec.Command("yt-dlp", "-o", "%(title)s.%(ext)s", "--quiet", url)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error executing yt-dlp: %w", err)
	}

	// get file name
	cmd = exec.Command("yt-dlp", "-o", "%(title)s.%(ext)s", "--print", "filename", "--no-warnings", url)
	stdOut, _ := cmd.CombinedOutput()
	filename := string(stdOut)
	filename = strings.Trim(filename, "\n")
	ext := filepath.Ext(filename)
	newFilename, err := utils.GenerateHash(filename)
	if err != nil {
		return "", fmt.Errorf("error generating hash: %w", err)
	}
	newFilename = fmt.Sprintf("%s%s", newFilename, ext)
	filepath_ := filepath.Join(outDir, newFilename)
	err = os.Rename(filename, filepath_)
	if err != nil {
		return "", err
	}

	return newFilename, nil
}
