package handler

import (
	"net/http"
	"strings"
)

func DenyDirectoryListing(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
}
