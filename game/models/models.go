package models

import (
	"time"

	"github.com/bpoetzschke/bin.go/game/word_manager"
)

type Game struct {
	ID             string
	RemainingWords []word_manager.Word
	FoundWords     []word_manager.FoundWord
	StartedAt      time.Time
	FinishedAt     *time.Time
}

type GameCollection []Game

func (gc GameCollection) GetCurrentGame() *Game {
	for _, g := range gc {
		if g.FinishedAt == nil {
			return &g
		}
	}

	return nil
}
