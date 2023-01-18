package mocks

import (
	"time"

	"github.com/NPeykov/snippetbox/internal/models"
)

var mockSnippet = &models.Snippet{
	ID: 0,
	Title: "The art of Mocking",
	Content: "The art of saying art",
	Created: time.Now(),
	Expires: time.Now().Add(2),
}

type SnippetModel struct {}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
    return 2, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
    switch id {
        case 1: return mockSnippet, nil
        default: return nil, models.ErrNoRecord
    }
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
    return []*models.Snippet{mockSnippet}, nil
}
