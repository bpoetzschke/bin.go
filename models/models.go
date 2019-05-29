package models

import (
	"strings"
	"time"
)

type Word struct {
	Value   string
	AddedBy string
	GifUrl  string
}

type WordList []Word

func (wl WordList) Join(sep string) string {
	words := []string{}

	for _, w := range wl {
		words = append(words, w.Value)
	}

	return strings.Join(words, sep)
}

type FoundWord struct {
	Word
	FoundBy string
}

type Game struct {
	ID             string
	RemainingWords []Word
	FoundWords     []FoundWord
	StartedAt      time.Time
	FinishedAt     *time.Time
}
