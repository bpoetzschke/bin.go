package main

import (
	"fmt"

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
	_, err := parseConfig()
	if err != nil {
		fmt.Printf("%s", err)
	}
}
