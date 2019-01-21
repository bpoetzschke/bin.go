package main

import (
	"fmt"
	"os"

	"github.com/bpoetzschke/bin.go/logger"

	"github.com/bpoetzschke/bin.go/game/word_manager"

	"github.com/bpoetzschke/bin.go/game"

	smw "github.com/bpoetzschke/bin.go/slack-middleware"
	"github.com/kelseyhightower/envconfig"
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

	word_manager.NewWordManager()

	game.NewGameLoop(mw.Connect()).Run()
}
