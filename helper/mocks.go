package helper

import (
	"github.com/bpoetzschke/bin.go/models"
	"github.com/stretchr/testify/mock"
)

type WordManagerMock struct {
	mock.Mock
}

func (m *WordManagerMock) LoadInitialWords() []models.Word {
	return m.Called().Get(0).([]models.Word)
}

func (m *WordManagerMock) GetWord(rawWord string) (models.Word, error) {
	called := m.Called(rawWord)
	return called.Get(0).(models.Word), called.Error(1)
}

func (m *WordManagerMock) GetGifForWord(word string) (string, error) {
	called := m.Called(word)
	return called.String(0), called.Error(1)
}
