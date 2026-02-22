package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as development progresses.
type application struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize a new instance of our application struct, containing the
	// dependencies (for now, just the structured logger).
	app := &application{
		logger: logger,
	}

	// serve static files such as css, js, and images

	logger.Info("starting server", "addr", *addr)
	// starts a web server. If this returns an err we use the log.Fatal() function to log the
	// error message and terminate the program.
	// Note: any error returned by http.ListenAndServe() is always non-nil
	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
