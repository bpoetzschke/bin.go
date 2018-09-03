package word_manager

import (
	"bufio"
	"fmt"
	"os"
	"sync"

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
	words        []word
}

func (wm *wordManager) init() {
	wm.loadWords()
}

func (wm *wordManager) loadWords() {
	words, err := wm.loadInitial()
	if err != nil {
		fmt.Printf("Error while loading words: %s", err)
	}

	chunkSize := (len(words) / concurrency) + 1
	fmt.Printf("Chunk size: %d\n", chunkSize)

	chunkIndex := 0

	wordMutex := sync.Mutex{}

	for chunkStart := 0; chunkStart < len(words); {

		end := chunkStart + chunkSize - 1

		if end > len(words)-1 {
			end = len(words) - 1
		}

		go func(chunk int, start int, end int) {
			//fmt.Printf("Handling chunk %d start %d end %d\n", chunk, start, end)
			for i := start; i <= end; i++ {
				//fmt.Printf("Process index: %d\n", i)
				url, err := wm.gifGenerator.Random(words[i])
				if err != nil {
					fmt.Printf("Error while fetching gif for word %q: %s\n", words[i], err)
					return
				}

				wordMutex.Lock()
				wm.words = append(wm.words, word{
					Word:   words[i],
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
