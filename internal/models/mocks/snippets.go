package mocks

import (
	"github.com/noonacedia/sourcepaste/internal/models"
	"time"
)

var MockSnippet = &models.Snippet{
	ID:      1,
	Title:   "TestSnippet",
	Content: "TestSnippet",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case MockSnippet.ID:
		return MockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{MockSnippet}, nil
}
