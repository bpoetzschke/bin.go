package models

import (
	"time"
)

type Word struct {
	Value   string
	AddedBy string
	GifUrl  string
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
