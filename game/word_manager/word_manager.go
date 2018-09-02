package word_manager

import (
	"bufio"
	"fmt"
	"os"
)

const initialWordFile = "initial.txt"

type WordManager interface {
}

func NewWordManager() WordManager {
	wm := wordManager{}
	wm.init()

	return &wm
}

type wordManager struct {
}

func (wm *wordManager) init() {
	wm.loadWords()
}

func (wm *wordManager) loadWords() {
	if err := wm.loadInitial(); err != nil {
		fmt.Printf("Error while loading words: %s", err)
	}
}

func (wm *wordManager) loadInitial() error {
	file, err := os.Open(initialWordFile)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var words []string

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return nil
}
