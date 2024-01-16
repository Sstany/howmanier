package core

import (
	"context"
	"database/sql"
	"errors"
	"howmanier/app/db"
	"howmanier/app/sdk"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	dbClient *db.PostgresClient
	bot      *tgbotapi.BotAPI
}

func (r *TelegramBot) getTgUserFromUpdate(update *tgbotapi.Update) *db.User {
	return &db.User{
		ID:       update.SentFrom().ID,
		Username: update.SentFrom().UserName,
	}
}

func NewTelegramBot(dbClient *db.PostgresClient) *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(os.Getenv(sdk.EnvToken))

	if err != nil {
		panic(err)
	}

	return &TelegramBot{
		dbClient: dbClient,
		bot:      bot,
	}
}
func (r *TelegramBot) process(update *tgbotapi.Update) error {
	tgUser := r.getTgUserFromUpdate(update)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := r.dbClient.FetchUser(ctx, tgUser); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = r.dbClient.CreateUser(ctx, tgUser)
		}

		if err != nil {
			return err
		}
	}

	switch update.Message.Command() {
	case "start":
		r.bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "hi"))
	}

	return nil
}

func (r *TelegramBot) Start() {
	config := tgbotapi.NewUpdate(0)
	config.Timeout = 60

	updates := r.bot.GetUpdatesChan(config)

	for update := range updates {
		r.process(&update)
	}

}
