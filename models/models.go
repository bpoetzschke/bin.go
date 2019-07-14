package models

import (
	"strings"
	"time"
)

type Word struct {
	Value   string `json:"value"`
	AddedBy string `json:"added_by"`
	GifUrl  string `json:"gif_url"`
}

type WordList []Word

func (wl WordList) Join(sep string) string {
	words := []string{}

	for _, w := range wl {
		words = append(words, w.Value)
	}

	return strings.Join(words, sep)
}

func (wl WordList) Diff(other WordList) WordList {
	s1, s2 := wl, other

	diff := WordList{}

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, a := range s1 {
			found := false
			for _, b := range s2 {
				if a.Value == b.Value {
					found = true
					break
				}
			}

			if !found {
				diff = append(diff, a)
			}
		}

		if i == 0 {
			s1, s2 = s2, s1
		}
	}

	return diff
}

type Game struct {
	ID             string     `json:"id"`
	RemainingWords WordList   `json:"remaining_words"`
	FoundWords     WordList   `json:"found_words"`
	StartedAt      time.Time  `json:"started_at"`
	FinishedAt     *time.Time `json:"finished_at"`
}

func (g *Game) AddNewWord(newWord Word) {
	g.RemainingWords = append(g.RemainingWords, newWord)
}
