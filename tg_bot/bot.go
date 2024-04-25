package tg_bot

import (
	"context"
	"fmt"
	"github.com/dunkbing/tinyimg/converter/config"
	"github.com/dunkbing/tinyimg/converter/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
	"os"
	"os/signal"
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

const (
	videoCommand           = "video"
	videoCommandPattern    = "/video"
	playlistCommand        = "playlist"
	playlistCommandPattern = "/playlist"
)

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
				Command:     videoCommand,
				Description: "Download video with the url",
			},
			{
				Command:     playlistCommand,
				Description: "Download a playlist with the url",
			},
		},
	})
	if err != nil {
		return
	}

	jobChan := make(chan Job, 100)
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypeExact,
		h.helpHandler,
	)
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		videoCommandPattern,
		bot.MatchTypePrefix,
		h.downloadHandler(videoCommand, jobChan),
	)
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		playlistCommandPattern,
		bot.MatchTypePrefix,
		h.downloadHandler(playlistCommand, jobChan),
	)

	go h.worker(b, jobChan)
	b.Start(ctx)
}

type Job struct {
	url     string
	chatID  int64
	command string
}

func (h *handler) worker(b *bot.Bot, jobChan <-chan Job) {
	slog.Info("Start processing job")
	for job := range jobChan {
		slog.Info("Processing job", "chat_id", job.chatID, "command", job.command, "url", job.url)
		var filename string
		var err error
		switch job.command {
		case videoCommand:
			filename, err = utils.DownloadVideo(job.url, h.config.App.OutDir)
		case playlistCommand:
			filename, err = utils.DownloadPlaylist(job.url, h.config.App.OutDir)
		}
		ctx := context.Background()
		if err != nil {
			slog.Error(fmt.Sprintf("Error downloading %s", job.command), "error", err)
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: job.chatID,
				Text:   "Something went wrong",
			})
			continue
		}
		slog.Info("Downloaded file", "filename", filename)

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: job.chatID,
			Text:   fmt.Sprintf("%s/video?f=%s", config.HostUrl, filename),
		})

		if err != nil {
			slog.Error("send message", slog.String("error", err.Error()))
		}
	}
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

func (h *handler) downloadHandler(command string, jobChan chan Job) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		url_ := update.Message.Text
		switch command {
		case videoCommand:
			if strings.HasPrefix(url_, videoCommandPattern) {
				url_ = strings.Replace(url_, videoCommandPattern, "", 1)
			}
		case playlistCommand:
			if strings.HasPrefix(url_, playlistCommandPattern) {
				url_ = strings.Replace(url_, playlistCommandPattern, "", 1)
			}
		default:
			slog.Error("Unknown command", "command", command)
			return
		}
		url_ = strings.TrimSpace(url_)
		if url_ == "" {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text: `Please use the following command pattern:
/video https://example.com/video.mp4
or
/playlist https://example.com/playlist`,
			})
			return
		}

		if !utils.IsValidUrl(url_) {
			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Invalid URL",
			})
			return
		}

		jobChan <- Job{
			url:     url_,
			chatID:  update.Message.Chat.ID,
			command: command,
		}
	}
}
