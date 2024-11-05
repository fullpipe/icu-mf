package main

import (
	"fmt"
	"htmx-templ/tr"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/text/language"
)

type GlobalState struct {
	Count int
}

var global GlobalState
var sessionManager *scs.SessionManager

func main() {
	// Initialize the session.
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("lang") {
			langTag, err := language.Parse(r.URL.Query().Get("lang"))
			if err != nil {
				slog.Error(err.Error(), slog.Any("lang", r.URL.Query().Get("lang")))
				langTag = language.English
			}

			sessionManager.Put(r.Context(), "lang", langTag.String())

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		userCount := sessionManager.GetInt(r.Context(), "count")
		component := page_apple_count(global.Count, userCount)

		component.Render(r.Context(), w)
	})

	mux.HandleFunc("/counters", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()

			if r.Form.Has("global") {
				global.Count++
			}

			if r.Form.Has("user") {
				currentCount := sessionManager.GetInt(r.Context(), "count")
				sessionManager.Put(r.Context(), "count", currentCount+1)
			}
		}

		userCount := sessionManager.GetInt(r.Context(), "count")
		component := counters_widget(global.Count, userCount)

		component.Render(r.Context(), w)
	})

	// Add the middleware.
	server := tr.LangMiddleware(sessionManager, mux)
	server = sessionManager.LoadAndSave(server)

	// Start the server.
	fmt.Println("listening on http://localhost:8000")
	if err := http.ListenAndServe(":8000", server); err != nil {
		log.Printf("error listening: %v", err)
	}
}
