package botkit

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	cmdViews map[string]ViewFunc
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{api: api}
}

func (b *Bot) Run(ctx context.Context) error {
	log.Print("[INFO] bot start")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(context.Background(), 5*time.Minute)
			b.handleUpdate(updateCtx, update)
			updateCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.cmdViews == nil {
		b.cmdViews = make(map[string]ViewFunc)
	}

	b.cmdViews[cmd] = view
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()

	var view ViewFunc

	if !update.Message.IsCommand() {
		return
	}

	cmd := update.Message.Command()

	cmdView, ok := b.cmdViews[cmd]
	if !ok {
		return
	}

	view = cmdView

	if err := view(ctx, b.api, update); err != nil {
		log.Printf("[ERROR] failed to handle update: %v", err)

		if _, err := b.api.Send(
			tgbotapi.NewMessage(update.Message.Chat.ID, "internal error"),
		); err != nil {
			log.Printf("[ERROR] failed to send message: %v", err)
		}
	}
}
