package mock

import (
	"time"

	"github.com/lipandr/Snippetbox/pkg/models"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "Hello World",
	Content: "This is a snippet",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (s *SnippetModel) Insert(_, _, _ string) (int, error) {
	return 2, nil
}

func (s *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (s *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
