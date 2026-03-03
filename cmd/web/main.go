package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/chrisdiebold/snippetbox/internal/dbx"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as development progresses.
type application struct {
	logger         *slog.Logger
	queries        *dbx.Queries
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	// standard sql.DB just for scs
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ping will fail fast if we cannot hit the database. usually caused by bad credentials
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("failed to reach database:", err)
	}

	queries := dbx.New(pool)

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	// Use the scs.New() function to initialize a new session manager. Then we
	// configure it to use our MySQL database as the session store, and set a
	// lifetime of 12 hours (so that sessions automatically expire 12 hours
	// after first being created).
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// Initialize a new instance of our application struct, containing the
	// dependencies (for now, just the structured logger).
	app := &application{
		logger:        logger,
		queries:       queries,
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		// Create a *log.Logger from our structured logger handler, which writes
		// log entries at the Error level, and assign it to the ErrorLog field. If
		// you would prefer to log the server errors at Warn level instead, you
		// could pass slog.LevelWarn as the final parameter.
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	// serve static files such as css, js, and images

	logger.Info("starting server", "addr", *addr)
	// starts a web server. If this returns an err we use the log.Fatal() function to log the
	// error message and terminate the program.
	// Note: any error returned by http.ListenAndServe() is always non-nil
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
