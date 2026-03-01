package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/chrisdiebold/snippetbox/internal/db"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as development progresses.
type application struct {
	logger        *slog.Logger
	queries       *db.Queries
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	// TODO: make this so it does not leak default values
	user := flag.String("user", "developer", "Database user name")
	password := flag.String("password", "developer", "Database password")

	flag.Parse()

	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	connStr := fmt.Sprintf("postgres://%s:%s@localhost:5432/snippetbox_dev?sslmode=disable", *user, *password)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal("failed to create pool:", err)
	}
	defer pool.Close()

	// ping will fail fast if we cannot hit the database. usually caused by bad credentials
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("failed to reach database:", err)
	}

	queries := db.New(pool)

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()
	// Initialize a new instance of our application struct, containing the
	// dependencies (for now, just the structured logger).
	app := &application{
		logger:        logger,
		queries:       queries,
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	// serve static files such as css, js, and images

	logger.Info("starting server", "addr", *addr)
	// starts a web server. If this returns an err we use the log.Fatal() function to log the
	// error message and terminate the program.
	// Note: any error returned by http.ListenAndServe() is always non-nil
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
