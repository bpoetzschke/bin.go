package word_manager

import (
	"bufio"
	"os"
	"sync"

	"github.com/bpoetzschke/bin.go/logger"
	"github.com/bpoetzschke/bin.go/models"

	"github.com/bpoetzschke/bin.go/gif"
)

const (
	initialWordFile = "initial.txt"
	concurrency     = 5
)

type WordManager interface {
	LoadWords() ([]models.Word, error)
}

func NewWordManager() WordManager {
	wm := wordManager{
		gifGenerator: gif.NewGiphy(),
	}

	return &wm
}

type wordManager struct {
	gifGenerator gif.Gif
}

func (wm *wordManager) LoadWords() ([]models.Word, error) {
	wordList, err := wm.loadInitial()
	if err != nil {
		logger.Error("Error while loading words: %s", err)
	}

	chunkSize := (len(wordList) / concurrency) + 1
	logger.Debug("Chunk size: %d", chunkSize)

	words := make([]models.Word, 0)

	chunkIndex := 0

	wordMutex := sync.Mutex{}
	waitGroup := sync.WaitGroup{}

	for chunkStart := 0; chunkStart < len(words); {

		end := chunkStart + chunkSize - 1

		if end > len(words)-1 {
			end = len(words) - 1
		}

		waitGroup.Add(1)

		go func(chunk int, start int, end int) {
			for i := start; i <= end; i++ {
				url, _, err := wm.gifGenerator.Random(wordList[i])
				if err != nil {
					logger.Warning("Could not fetch gif for word %q: %s", words[i], err)
					continue
				}

				wordMutex.Lock()
				words = append(words, models.Word{
					Value:  wordList[i],
					GifUrl: url,
				})
				wordMutex.Unlock()

			}

			waitGroup.Done()
		}(chunkIndex, chunkStart, end)

		// update chunk start for next iteration of the loop
		chunkStart = end + 1
		chunkIndex++
	}

	waitGroup.Wait()
	return words, nil
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
