package core

import (
	"context"
	"database/sql"
	"errors"
	"howmanier/app/db"
	"howmanier/app/sdk"
	"math/rand"

	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type TelegramBot struct {
	dbClient *db.PostgresClient
	logger   *zap.Logger
	bot      *tgbotapi.BotAPI
}

func (r *TelegramBot) getTgUserFromUpdate(update *tgbotapi.Update) *db.User {
	return &db.User{
		ID:       update.SentFrom().ID,
		Username: update.SentFrom().UserName,
	}
}

func NewTelegramBot(dbClient *db.PostgresClient, logger *zap.Logger) *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(os.Getenv(sdk.EnvToken))
	if err != nil {
		panic(err)
	}

	return &TelegramBot{
		dbClient: dbClient,
		bot:      bot,
		logger:   logger,
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
		r.bot.Send(tgbotapi.NewMessage(update.Message.From.ID,"Привет, я бот для контроля продуктов в твоём холодильнике
		/add- добавит продукт
		/delite - удалить продукт из списка
		/list - показать актуальный список продуктов
		/whattoeat - поможет выбрать, что приготовить сегодня"))
	case "add":
		r.handlerAdd(ctx, tgUser, update)
	case "delite":
		r.handlerDelite(ctx, tgUser, update)
	case "list":
		r.handlerList(ctx, tgUser, update)
	case "whattoeat":
		r.handlerWhattoeat(update)
	}

	return nil
}
func (r *TelegramBot) handlerAdd(ctx context.Context, user *db.User, update *tgbotapi.Update) {
	name := update.Message.CommandArguments()
	r.bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "Добавил "+name))

}

func (r *TelegramBot) handlerDelite(ctx context.Context, user *db.User, update *tgbotapi.Update) {
	name := update.Message.CommandArguments()
	r.bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "Удалил "+name))
}

func (r *TelegramBot) handlerList(ctx context.Context, user *db.User, update *tgbotapi.Update) {
	//name:= update.Message.CommandArguments()
	r.bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "Твой список: "))
}

func (r *TelegramBot) handlerWhattoeat(update *tgbotapi.Update) {
	dish := []string{"пельмени", "шарлотка", "рис с мяосм", "голубцы", "суп с тефтельками", "курица с овощами"}
	r.bot.Send(tgbotapi.NewMessage(update.Message.From.ID, dish[rand.Intn(6)]))

}
func (r *TelegramBot) Start() {
	var err error

	config := tgbotapi.NewUpdate(0)
	config.Timeout = 60

	updates := r.bot.GetUpdatesChan(config)

	for update := range updates {
		err = r.process(&update)
		if err != nil {
			r.logger.Error("failed to process update", zap.Error(err))
		}
	}

}
