package storage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bpoetzschke/bin.go/models"

	"github.com/stretchr/testify/require"
)

func TestFileStorageLoadFileWithWrongJson(t *testing.T) {
	require.NoError(t, ioutil.WriteFile(fileName, []byte(`I am not a json.`), 0644))
	defer os.Remove(fileName)

	s := fileStorage{}
	content := s.loadFile()
	require.EqualValues(t, fileModel{
		GameMap:  map[string]models.Game{},
		WordList: []models.Word{},
	}, content)
}
