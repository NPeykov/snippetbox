package main

import "github.com/NPeykov/snippetbox/internal/models"

type templateData struct {
    Snippet *models.Snippet
    Snippets []*models.Snippet
}
