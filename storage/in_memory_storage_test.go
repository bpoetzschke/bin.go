package storage_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"

	"github.com/bpoetzschke/bin.go/models"
	"github.com/bpoetzschke/bin.go/storage"
)

func TestInMemoryStorage(t *testing.T) {
	ims := storage.NewInMemoryStorage()

	// when storage is initialized the word list should be empty
	words, err := ims.LoadWordList()
	require.NoError(t, err)
	require.EqualValues(t, []models.Word{}, words)

	// adding a word
	word1 := models.Word{Value: uuid.NewV4().String()}
	success, err := ims.AddWord(word1)
	require.NoError(t, err)
	require.True(t, success)

	// add same word again, this should fail
	success, err = ims.AddWord(word1)
	require.Error(t, err)
	require.False(t, success)

	//retrieve word list and check if word exists
	words, err = ims.LoadWordList()
	require.NoError(t, err)
	require.EqualValues(t, []models.Word{word1}, words)
}
