package storage

import (
	"github.com/bpoetzschke/bin.go/models"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (m *StorageMock) LoadWordList() ([]models.Word, error) {
	called := m.Called()

	return called.Get(0).([]models.Word), called.Error(1)
}

func (m *StorageMock) AddWord(word models.Word) (bool, models.Word, error) {
	called := m.Called(word)

	return called.Bool(0), called.Get(1).(models.Word), called.Error(2)
}

func (m *StorageMock) LoadCurrentGame() (models.Game, bool, error) {
	called := m.Called()

	return called.Get(0).(models.Game), called.Bool(1), called.Error(2)
}

func (m *StorageMock) SaveGame(game models.Game) error {
	return m.Called(game).Error(0)
}
