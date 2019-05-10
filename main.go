package main

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"

	"github.com/bpoetzschke/bin.go/game"
	"github.com/bpoetzschke/bin.go/logger"
	smw "github.com/bpoetzschke/bin.go/slack-middleware"
	"github.com/bpoetzschke/bin.go/storage"
	"github.com/bpoetzschke/bin.go/word_manager"
)

type config struct {
	SlackToken string `split_words:"true"`
	Debug      bool   `default:"false"`
}

func parseConfig() (config, error) {

	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return config{}, fmt.Errorf("Failed to parse environment config.\n%s\n", err)
	}

	if len(cfg.SlackToken) == 0 {
		return config{}, fmt.Errorf("required property %q is not set\n", "SLACK_TOKEN")
	}

	return cfg, nil
}

func setDebug() {
	logger.Enable()
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	mw := smw.NewMiddleware(cfg.SlackToken)

	if cfg.Debug {
		setDebug()
	}

	game, err := game.NewGameLoop(
		mw,
		storage.NewInMemoryStorage(),
		word_manager.NewWordManager(),
	)
	if err != nil {
		logger.Error("Error while starting game. %s", err)
		return
	}

	game.Run()
}
