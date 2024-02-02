package core

import (
	"context"
	"database/sql"
	"errors"
	"howmanier/app/db"
	"howmanier/app/sdk"
	"math/rand"
	"strconv"
	"strings"

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
		r.bot.Send(tgbotapi.NewMessage(update.Message.From.ID, `Привет, я бот для контроля продуктов в твоём холодильнике`))
	case "add":
		return r.handlerAdd(ctx, tgUser, update)
	case "delete":
		r.handlerDelete(ctx, tgUser, update)
	case "list":
		return r.handlerList(ctx, tgUser, update)
	case "whattoeat":
		r.handlerWhattoeat(update)
	}

	return nil
}
func (r *TelegramBot) handlerAdd(ctx context.Context, user *db.User, update *tgbotapi.Update) error {
	args := update.Message.CommandArguments()
	switch {
	case strings.Contains(args, " "):
		food := strings.Split(args, " ")

		count, err := strconv.Atoi(food[1])
		if err != nil {
			return err
		}

		if err := r.dbClient.AddProduct(ctx,
			&db.Product{
				UserID: user.ID,
				Name:   food[0],
				Count:  count,
			},
		); err != nil {
			return err
		}

		r.bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "Добавил "+food[0]+" "+food[1]+"шт."))

	case strings.Contains(args, " ") == false:
		r.bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "Добавьте количество продукта или сам продукт"))

	}

	return nil
}

func (r *TelegramBot) handlerDelete(ctx context.Context, user *db.User, update *tgbotapi.Update) error {
	args := update.Message.CommandArguments()

	switch {
	case strings.Contains(args, " "):
		food := strings.Split(args, " ")

		count, err := strconv.Atoi(food[1])
		if err != nil {
			return err
		}

		if err := r.dbClient.DeleteProduct(ctx,
			&db.Product{
				UserID: user.ID,
				Name:   food[0],
				Count:  count,
			},
		); err != nil {
			return err
		}

		r.bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "Удалил "+food[0]+" "+food[1]+"шт."))

	case strings.Contains(args, " ") == false:
		r.bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "Добавьте количество продукта или сам продукт"))

	}

	return nil
}

func (r *TelegramBot) handlerList(ctx context.Context, user *db.User, update *tgbotapi.Update) error {
	//name:= update.Message.CommandArguments()
	products, err := r.dbClient.ListFridge(ctx, user)
	if err != nil {
		r.bot.Send(tgbotapi.NewMessage(update.Message.From.ID, "АШИПКА"))
		return err
	}

	out := strings.Builder{}
	for i, product := range products {
		out.WriteString(strconv.Itoa(i + 1))
		out.WriteString(".  ")
		out.WriteString(product.Name)
		out.WriteString(" ")
		out.WriteString(strconv.Itoa(product.Count))
		out.WriteString("\n")

	}

	r.bot.Send(tgbotapi.NewMessage(update.Message.From.ID, out.String()))

	return nil
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
