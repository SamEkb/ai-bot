package main

import (
	"app/internal/config"
	"app/internal/controller/bot"
	"app/internal/service/openai"
	"context"
	"github.com/theartofdevel/logging"
	"os"
)

func main() {
	ctx := context.Background()

	cfg := config.GetConfig()

	logging.Default().Info(
		"Application starts with configuration",
		logging.StringAttr("app_id", cfg.App.Id),
		logging.StringAttr("app_name", cfg.App.Name),
	)

	logger := logging.NewLogger(
		logging.WithLevel(cfg.App.LogLevel),
		logging.WithIsJSON(cfg.App.IsLogJSON),
	)

	ctx = logging.ContextWithLogger(ctx, logger)

	openaiCfg := openai.NewConfig(cfg.OpenAIConfig.ApiKey)

	openaiService, err := openai.NewService(openaiCfg)
	if err != nil {
		logging.WithAttrs(ctx, logging.ErrAttr(err)).Error("failed to create openai service")
		os.Exit(1)
	}

	botConfig := bot.NewConfig(cfg.Bot.Token, cfg.Bot.Timeout)
	botWrapper, err := bot.NewWrapper(botConfig, openaiService)
	if err != nil {
		logging.WithAttrs(ctx, logging.ErrAttr(err)).Error("failed to create bot wrapper")
		os.Exit(1)
	}

	err = botWrapper.Start(ctx)
	if err != nil {
		logging.WithAttrs(ctx, logging.ErrAttr(err)).Error("failed to start bot")
		os.Exit(1)
	}
}
