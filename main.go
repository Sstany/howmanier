package main

import (
	"context"
	"os"

	"howmanier/app/core"
	"howmanier/app/db"
	"howmanier/app/sdk"
)

func main() {
	postgresClient := db.NewPostgresClient(context.Background(), os.Getenv(sdk.EnvPostgres))
	bot := core.NewTelegramBot(postgresClient)
	bot.Start()

}
