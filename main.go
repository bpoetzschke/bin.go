package main

import (
	"fmt"
	"os"

	smw "github.com/bpoetzschke/bin.go/slack-middleware"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	SlackToken string `split_words:"true"`
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

func main() {
	cfg, err := parseConfig()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	mw := smw.NewMiddleware(cfg.SlackToken)
	events := mw.Connect()
	for evt := range events{
		fmt.Printf("Received event: %+v", evt)
	}

	fmt.Printf("")
}
