package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/render"

	"github.com/jonnylangefeld/go-api/pkg/types"
)

const (
	// PageIDKey refers to the context key that stores the next page id
	PageIDKey CustomKey = "page_id"
)

// Pagination middleware is used to extract the next page id from the url query
func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		PageID := r.URL.Query().Get(string(PageIDKey))
		intPageID := 0
		var err error
		if PageID != "" {
			intPageID, err = strconv.Atoi(PageID)
			if err != nil {
				_ = render.Render(w, r, types.ErrInvalidRequest(fmt.Errorf("couldn't read %s: %w", PageIDKey, err)))
				return
			}
		}
		ctx := context.WithValue(r.Context(), PageIDKey, intPageID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
