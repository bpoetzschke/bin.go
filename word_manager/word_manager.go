package word_manager

import (
	"bufio"
	"os"
	"sync"
	"time"

	"github.com/bpoetzschke/bin.go/logger"
	"github.com/bpoetzschke/bin.go/models"

	"github.com/bpoetzschke/bin.go/gif"
)

const (
	initialWordFile = "initial.txt"
	concurrency     = 2
)

type WordManager interface {
	LoadInitialWords() []models.Word
	GetGifForWord(word string) (string, error)
}

func NewWordManager() (WordManager, error) {
	gifGenerator, err := gif.NewGiphy()
	if err != nil {
		return nil, err
	}

	wm := wordManager{
		gifGenerator: gifGenerator,
	}

	return &wm, nil
}

type wordManager struct {
	gifGenerator gif.Gif
}

func (wm *wordManager) GetGifForWord(word string) (string, error) {
	url, _, err := wm.gifGenerator.Random(word)

	return url, err
}

func (wm *wordManager) LoadInitialWords() []models.Word {
	wordList, err := wm.loadFromInitialFile()
	if err != nil {
		logger.Error("Error while loading words: %s", err)
	}

	chunkSize := (len(wordList) / concurrency) + 1
	logger.Debug("Chunk size: %d", chunkSize)

	words := make([]models.Word, 0)

	chunkIndex := 0

	wordMutex := sync.Mutex{}
	waitGroup := sync.WaitGroup{}

	chunkStart := 0

	for chunkStart < len(wordList) {

		end := chunkStart + chunkSize - 1

		if end > len(wordList)-1 {
			end = len(wordList) - 1
		}

		waitGroup.Add(1)

		go func(chunk int, start int, end int) {
			for i := start; i <= end; i++ {
				url, found, err := wm.gifGenerator.Random(wordList[i])
				if err != nil {
					logger.Warning("Could not fetch gif for word %q: %s", wordList[i], err)
					continue
				}

				if !found {
					logger.Info("Could not find gif for word %q.", wordList[i])
				} else {
					logger.Debug("Loaded word: %q, url %q", wordList[i], url)
				}

				wordMutex.Lock()
				words = append(words, models.Word{
					Value:  wordList[i],
					GifUrl: url,
				})
				wordMutex.Unlock()
				<-time.After(100 * time.Millisecond)
			}

			waitGroup.Done()
		}(chunkIndex, chunkStart, end)

		// update chunk start for next iteration of the loop
		chunkStart = end + 1
		chunkIndex++
	}

	waitGroup.Wait()
	return words
}

func (wm *wordManager) loadFromInitialFile() ([]string, error) {
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
