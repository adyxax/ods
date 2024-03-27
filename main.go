package main

import (
	"embed"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"unicode"
)

// Variables to customise the search behaviour
const (
	listenStr = "0.0.0.0:8090"
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
	HasQuery bool
	Query    string
	Invalid  bool
}

func getIndex() http.Handler {
	data := IndexData{
		HasQuery: false,
		Query:    "",
		Invalid:  true,
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
			HasQuery: true,
			Query:    r.FormValue("query"),
			Invalid:  true,
		}
		query, _, _ := transform.String(normalizer, strings.TrimSpace(strings.ToUpper(data.Query)))
		for _, w := range words {
			if w == query {
				data.Invalid = false
				break
			}
		}
		w.Header().Set("Cache-Control", "no-store, no-cache")
		indexTemplate.ExecuteTemplate(w, "index.html", data)
	})
}

// The main function
func main() {
	http.Handle("GET /static/", http.FileServer(http.FS(static)))
	http.Handle("GET /", getIndex())
	http.Handle("POST /", postIndex())
	slog.Info("listening", "addr", listenStr)
	if err := http.ListenAndServe(listenStr, nil); err != nil && err != http.ErrServerClosed {
		slog.Error("error listening and serving", "error", err)
	}
}
