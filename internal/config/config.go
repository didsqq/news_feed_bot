package config

import (
	"log"
	"sync"
	"time"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfighcl"
)

type Config struct {
	TelegramBotToken     string        `hcl:"telegram_bot_token" env:"TELEGRAM_BOT_TOKEN" required:"true"`
	TelegramChannelID    int64         `hcl:"telegram_channel_id" env:"TELEGRAM_CHANNEL_ID" required:"true"`
	DatabaseDSN          string        `hcl:"database_dsn" env:"DATABASE_DSN" default:"postgres://postgres:postgres@localhost:5432/news_feed_bot?sslmode=disable"`
	FetchInterval        time.Duration `hcl:"fetch_interval" env:"FETCH_INTERVAL" default:"10m"`
	NotificationInterval time.Duration `hcl:"notification_interval" env:"NOTIFICATION_INTERVAL" default:"1m"`
	FilterKeywords       []string      `hcl:"filter_keywords" env:"FILTER_kEYWORDS"`
	OpenAIKey            string        `hcl:"openai_key" env:"OPENAI_KEY"`
	OpenAIPromзt         string        `hcl:"openai_prompt" env:"OPENAI_PROMPT"`
	OpenAШModel          string        `hcl:"openai_model" env:"OPENAI_MODEL" default:"gpt-3.5-turbo"`
}

var (
	cfg  Config    // глобальная переменная типа Config, которая будет содержать загруженную конфигурацию
	once sync.Once // используется для того, чтобы гарантировать, что инициализация произойдет только один раз
)

func Get() Config {
	once.Do(func() {
		loader := aconfig.LoaderFor(&cfg, aconfig.Config{
			EnvPrefix: "NFB", // переменные окружения начинающиеся с префикса NFB_
			Files:     []string{"./config.hcl", "./config.local.hcl", "$HOME/.config/news-feed-bot/config.hcl"},
			FileDecoders: map[string]aconfig.FileDecoder{ // используется декодер aconfighcl
				".hcl": aconfighcl.New(),
			},
		})

		if err := loader.Load(); err != nil {
			log.Printf("[ERROR] failed to load config: %v", err)
		}
	})

	return cfg
}