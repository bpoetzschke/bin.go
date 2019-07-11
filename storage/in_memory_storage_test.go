package storage_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"

	"github.com/bpoetzschke/bin.go/models"
	"github.com/bpoetzschke/bin.go/storage"
)

func TestInMemoryWordStorage(t *testing.T) {
	ims := storage.NewInMemoryStorage()

	// when storage is initialized the word list should be empty
	words, err := ims.LoadWordList()
	require.NoError(t, err)
	require.EqualValues(t, []models.Word{}, words)

	// adding a word
	word1 := models.Word{Value: uuid.NewV4().String()}
	success, _, err := ims.AddWord(word1)
	require.NoError(t, err)
	require.True(t, success)

	// add same word again, this should fail
	success, existingWord, err := ims.AddWord(word1)
	require.NoError(t, err)
	require.False(t, success)
	require.EqualValues(t, word1, existingWord)

	//retrieve word list and check if word exists
	words, err = ims.LoadWordList()
	require.NoError(t, err)
	require.EqualValues(t, []models.Word{word1}, words)
}

func TestInMemoryGameStorage(t *testing.T) {
	ims := storage.NewInMemoryStorage()

	// when retrieving current we should get nil because there is no game stored
	game, found, err := ims.LoadCurrentGame()
	require.NoError(t, err)
	require.False(t, found)
	require.Empty(t, game)

	// create a game and store it
	game = models.Game{
		ID:        uuid.NewV4().String(),
		StartedAt: time.Now().UTC(),
	}
	err = ims.SaveGame(game)
	require.NoError(t, err)

	// retrieve current game
	value, found, err := ims.LoadCurrentGame()
	require.NoError(t, err)
	require.True(t, found)
	require.EqualValues(t, game, value)

	// update game and set it to finish and retrieve game afterwards --> current game should be empty because there is
	// no active game remaining
	now := time.Now().UTC()
	game.FinishedAt = &now
	err = ims.SaveGame(game)
	require.NoError(t, err)

	value, found, err = ims.LoadCurrentGame()
	require.NoError(t, err)
	require.False(t, found)
	require.Empty(t, value)
}
