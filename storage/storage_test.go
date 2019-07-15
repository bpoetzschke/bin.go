package storage

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/bpoetzschke/bin.go/models"
	"github.com/twinj/uuid"
)

type StorageTestSuite struct {
	suite.Suite
}

func (suite *StorageTestSuite) TearDownTest() {
	suite.Require().NoError(os.Remove(fileName))
}

func (suite *StorageTestSuite) TestWordStorage() {
	for storageMethod := range supportedStorageMethods {
		suite.T().Logf("Testing storage method: %s.", storageMethod)
		s := NewStorage(string(storageMethod))

		// when storage is initialized the word list should be empty
		words, err := s.LoadWordList()
		suite.Require().NoError(err)
		suite.Require().EqualValues([]models.Word{}, words)

		// adding a word
		word1 := models.Word{Value: uuid.NewV4().String()}
		success, _, err := s.AddWord(word1)
		suite.Require().NoError(err)
		suite.Require().True(success)

		// add same word again, this should fail
		success, existingWord, err := s.AddWord(word1)
		suite.Require().NoError(err)
		suite.Require().False(success)
		suite.Require().EqualValues(word1, existingWord)

		//retrieve word list and check if word exists
		words, err = s.LoadWordList()
		suite.Require().NoError(err)
		suite.Require().EqualValues([]models.Word{word1}, words)
	}

}

func (suite *StorageTestSuite) TestGameStorage() {
	for storageMethod := range supportedStorageMethods {
		suite.T().Logf("Testing storage method: %s.", storageMethod)
		s := NewStorage(string(storageMethod))

		// when retrieving current we should get nil because there is no game stored
		game, found, err := s.LoadCurrentGame()
		suite.Require().NoError(err)
		suite.Require().False(found)
		suite.Require().Empty(game)

		// create a game and store it
		game = models.Game{
			ID:        uuid.NewV4().String(),
			StartedAt: time.Now().UTC(),
		}
		err = s.SaveGame(game)
		suite.Require().NoError(err)

		// retrieve current game
		value, found, err := s.LoadCurrentGame()
		suite.Require().NoError(err)
		suite.Require().True(found)
		suite.Require().EqualValues(game, value)

		// update game and set it to finish and retrieve game afterwards --> current game should be empty because there is
		// no active game remaining
		now := time.Now().UTC()
		game.FinishedAt = &now
		err = s.SaveGame(game)
		suite.Require().NoError(err)

		value, found, err = s.LoadCurrentGame()
		suite.Require().NoError(err)
		suite.Require().False(found)
		suite.Require().Empty(value)
	}
}

func TestStorage(t *testing.T) {
	suite.Run(t, &StorageTestSuite{})
}
