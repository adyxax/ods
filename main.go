package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

//go:embed ods.txt
var ods string

//go:embed index.html
var templatesFS embed.FS

//go:embed static/*
var static embed.FS

// html templates
var indexTemplate = template.Must(template.New("index").ParseFS(templatesFS, "index.html"))

type IndexData struct {
	Invalid bool
	Query   string
}

func getIndex() http.Handler {
	data := IndexData{
		Query:   "",
		Invalid: true,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache")
		indexTemplate.ExecuteTemplate(w, "index.html", data)
	})
}

func postIndex() http.Handler {
	normalizer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	words := strings.Split(ods, "\n")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := IndexData{
			Query:   r.FormValue("query"),
			Invalid: true,
		}
		query, _, _ := transform.String(normalizer, strings.TrimSpace(strings.ToUpper(data.Query)))
		for _, w := range words {
			if w == query {
				data.Invalid = false
				break
			}
		}
		slog.Info("post", "word", query, "invalid", data.Invalid)
		indexTemplate.ExecuteTemplate(w, "index.html", data)
	})
}

func run(
	ctx context.Context,
	getenv func(string) string,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServer(http.FS(static)))
	mux.Handle("GET /", getIndex())
	mux.Handle("POST /", postIndex())

	host := getenv("ODS_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := getenv("ODS_PORT")
	if port == "" {
		port = "8080"
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}
	errChan := make(chan error, 1)
	go func() {
		slog.Info("backend http server listening", "address", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("error listening and serving backend http server on %s: %w", httpServer.Addr, err)
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		slog.Info("backend http server shutting down")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("error shutting down backend http server: %w", err)
		}
	}

	return nil
}

func main() {
	ctx := context.Background()

	var opts *slog.HandlerOptions
	if os.Getenv("ODS_DEBUG") != "" {
		opts = &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	if err := run(
		ctx,
		os.Getenv,
	); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
