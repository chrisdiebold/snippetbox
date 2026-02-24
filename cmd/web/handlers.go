package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/chrisdiebold/snippetbox/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	// ParseFiles reads the template file into a template set. Then,
	// Execute() writes the template content as the response body.
	// The last arg to Execute() represents dynamic data we want to pass in.
	w.Header().Add("Server", "Go")

	snippets, err := app.queries.GetActiveSnippetsLimit10(ctx)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, http.StatusOK, "home.tmpl.html", templateData{
		Snippets: snippets,
	})
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	// always validate untrusted user input
	// in this case, must be a positive integer
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.queries.GetSnippetNotExpired(ctx, int32(id))
	if err != nil {
		app.logger.Error(db.ErrNoRecord.Error())
		app.serverError(w, r, err)
		return
	}
	// Use the new render helper.
	app.render(w, r, http.StatusOK, "view.tmpl.html", templateData{
		Snippet: snippet,
	})

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from snippet create"))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	// Create some variables holding dummy data. We'll remove these later on
	// during development.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	expiration := pgtype.Timestamptz{
		Time:  time.Now().Add(time.Duration(expires) * 24 * time.Hour),
		Valid: true,
	}
	created := pgtype.Timestamptz{
		Time:  time.Now(),
		Valid: true,
	}

	params := db.CreateSnippetParams{
		Title:   title,
		Content: content,
		Created: created,
		Expires: expiration,
	}

	s, err := app.queries.CreateSnippet(ctx, params)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", s.ID), http.StatusSeeOther)
}
