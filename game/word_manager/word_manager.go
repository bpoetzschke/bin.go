package word_manager

import (
	"bufio"
	"os"
	"sync"

	"github.com/bpoetzschke/bin.go/logger"

	"github.com/bpoetzschke/bin.go/game/gif"
)

const (
	initialWordFile = "initial.txt"
	concurrency     = 5
)

type WordManager interface {
}

func NewWordManager() WordManager {
	wm := wordManager{
		gifGenerator: gif.NewGiphy(),
	}
	wm.init()

	return &wm
}

type wordManager struct {
	gifGenerator gif.Gif
	words        []Word
}

func (wm *wordManager) init() {
	wm.loadWords()
}

func (wm *wordManager) loadWords() {
	words, err := wm.loadInitial()
	if err != nil {
		logger.Error("Error while loading words: %s", err)
	}

	chunkSize := (len(words) / concurrency) + 1
	logger.Debug("Chunk size: %d", chunkSize)

	chunkIndex := 0

	wordMutex := sync.Mutex{}

	for chunkStart := 0; chunkStart < len(words); {

		end := chunkStart + chunkSize - 1

		if end > len(words)-1 {
			end = len(words) - 1
		}

		go func(chunk int, start int, end int) {
			for i := start; i <= end; i++ {
				url, err := wm.gifGenerator.Random(words[i])
				if err != nil {
					logger.Warning("Could not fetch gif for word %q: %s", words[i], err)
					continue
				}

				wordMutex.Lock()
				wm.words = append(wm.words, Word{
					Value:  words[i],
					GifUrl: url,
				})
				wordMutex.Unlock()

			}
		}(chunkIndex, chunkStart, end)

		// update chunk start for next iteration of the loop
		chunkStart = end + 1
		chunkIndex++
	}
}

func (wm *wordManager) loadInitial() ([]string, error) {
	file, err := os.Open(initialWordFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var words []string

	for scanner.Scan() {
		word := scanner.Text()
		words = append(words, word)
	}

	return words, nil
}
