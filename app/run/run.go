package run

import "go.uber.org/zap"

var Logger *zap.Logger

func Init() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	Logger = logger
}
