package main

import "github.com/chrisdiebold/snippetbox/internal/db"

type templateData struct {
	Snippet  db.Snippet
	Snippets []db.Snippet
}
