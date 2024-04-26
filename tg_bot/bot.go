package tg_bot

import (
	"context"
	"fmt"
	"github.com/dunkbing/tinyimg/tinyimg/cache"
	"github.com/dunkbing/tinyimg/tinyimg/config"
	"github.com/dunkbing/tinyimg/tinyimg/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

type handler struct {
	config *config.Config
	redis  *redis.Client
}

func New() *handler {
	c := config.GetConfig()
	r := cache.GetRedisClient()
	return &handler{
		config: c,
		redis:  r,
	}
}

const (
	botName                = "tg_video_downloader"
	videoCommand           = "video"
	videoCommandPattern    = "/video"
	playlistCommand        = "playlist"
	playlistCommandPattern = "/playlist"
	statsCommand           = "stats"
	statsCommandPattern    = "/stats"
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
				Command:     "help",
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
			{
				Command:     statsCommand,
				Description: "Get bot's statistic",
			},
		},
	})
	if err != nil {
		return
	}

	jobChan := make(chan Job, 100)
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/help",
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
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		statsCommandPattern,
		bot.MatchTypeExact,
		h.statsHandler,
	)

	go h.processJob(b, jobChan)
	b.Start(ctx)
}

type Job struct {
	url     string
	chatID  int64
	command string
}

func (h *handler) processJob(b *bot.Bot, jobChan <-chan Job) {
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
			continue
		}

		val, err := h.redis.Get(ctx, botName).Result()
		if err != nil {
			slog.Error("redis get", slog.String("error", err.Error()))
		}
		if val == "" {
			h.redis.Set(ctx, botName, 1, 0)
		} else {
			vidDownloaded, err := strconv.Atoi(val)
			if err != nil {
				slog.Error("redis get", slog.String("error", err.Error()))
			}
			vidDownloaded += 1
			h.redis.Set(ctx, botName, vidDownloaded, 0)
		}
	}
}

func (h *handler) helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	help := `
🤖 Video Downloader Bot 🔽

📥 /video <url> - Download a single video
📥 /playlist <url> - Download a playlist

Supported Sources:
YouTube, TikTok, Facebook, Pornhub

Example:
/video https://example.com/video.mp4
/playlist https://example.com/playlist

💬 For support:
TikTok: tiktok.com/@dunkbing
Facebook: fb.com/dunkbinggg
Instagram: instagram.com/dunkbingg
X: x.com/dunkbingg
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

func (h *handler) statsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	val, err := h.redis.Get(ctx, botName).Result()
	if err != nil {
		slog.Error("Stats get", slog.String("error", err.Error()))
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Something went wrong",
		})
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Total video downloaded: %s", val),
	})
}
