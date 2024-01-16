package main

import (
	"context"
	"os"

	"howmanier/app/core"
	"howmanier/app/db"
	"howmanier/app/run"
	"howmanier/app/sdk"
)

func main() {
	run.Init()

	postgresClient := db.NewPostgresClient(context.Background(), os.Getenv(sdk.EnvPostgres))
	bot := core.NewTelegramBot(postgresClient, run.Logger)
	bot.Start()

}
