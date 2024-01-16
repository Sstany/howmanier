package core

import (
	"howmanier/app/db"
	"howmanier/app/sdk"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	dbClient *db.PostgresClient
	bot      *tgbotapi.BotAPI
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
