// cmd/seed/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/chrisdiebold/snippetbox/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func main() {
	user := flag.String("user", "developer", "Database user name")
	password := flag.String("password", "developer", "Database password")

	flag.Parse()

	ctx := context.Background()

	connStr := fmt.Sprintf("postgres://%s:%s@localhost:5432/snippetbox_dev?sslmode=disable", *user, *password)

	conn, err := pgx.Connect(ctx, connStr)

	if err != nil {
		log.Fatal("failed to connect:", err)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	if err := seed(ctx, queries); err != nil {
		log.Fatal("seeding failed:", err)
	}

	log.Println("database seeded successfully")
}

func seed(ctx context.Context, queries *db.Queries) error {
	expires := pgtype.Timestamptz{
		Time:  time.Now().Add(365 * 24 * time.Hour),
		Valid: true,
	}
	now := pgtype.Timestamptz{
		Time:  time.Now(),
		Valid: true,
	}

	snippets := []db.CreateSnippetParams{
		{
			Title:   "An old silent pond",
			Content: "An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō",
			Created: now,
			Expires: expires,
		},

		{
			Title:   "Over the wintry forest",
			Content: "Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki",
			Created: now,
			Expires: expires,
		},
		{
			Title:   "First autumn morning",
			Content: "First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo",
			Created: now,
			Expires: expires,
		},
	}

	for _, u := range snippets {
		snippet, err := queries.CreateSnippet(ctx, u)
		if err != nil {
			return fmt.Errorf("creating snippet %s: %w", u.Title, err)
		}
		log.Printf("created snippet: %s (%d)\n", snippet.Title, snippet.ID)
	}

	return nil
}
