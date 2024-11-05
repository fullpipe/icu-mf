package tr

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

// LangMiddleware bypasses user language to request context
func LangMiddleware(sessionManager *scs.SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := sessionManager.GetString(r.Context(), "lang")
		r = r.WithContext(context.WithValue(r.Context(), langKey, lang))

		next.ServeHTTP(w, r)
	})
}
