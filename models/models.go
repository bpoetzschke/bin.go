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
	ID             string
	RemainingWords WordList
	FoundWords     WordList
	StartedAt      time.Time
	FinishedAt     *time.Time
}
