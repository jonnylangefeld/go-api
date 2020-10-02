package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/jonnylangefeld/go-api/pkg/db"
	"github.com/jonnylangefeld/go-api/pkg/types"
)

type (
	// CustomKey is used to refer to the context key that stores custom values of this api to avoid overwrites
	CustomKey string
)

const (
	// ArticleCtxKey refers to the context key that stores the article
	ArticleCtxKey CustomKey = "article"
	// OrderCtxKey refers to the context key that stores the order
	OrderCtxKey CustomKey = "order"
)

var DBClient db.ClientInterface

func SetDBClient(c db.ClientInterface) {
	DBClient = c
}

// Article middleware is used to load an Article object from
// the URL parameters passed through as the request. In case
// the Article could not be found, we stop here and return a 404.
func Article(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var article *types.Article

		if id := chi.URLParam(r, "id"); id != "" {
			intID, err := strconv.Atoi(id)
			if err != nil {
				_ = render.Render(w, r, types.ErrInvalidRequest(err))
				return
			}
			article = DBClient.GetArticleByID(intID)
		} else {
			_ = render.Render(w, r, types.ErrNotFound())
			return
		}
		if article == nil {
			_ = render.Render(w, r, types.ErrNotFound())
			return
		}

		ctx := context.WithValue(r.Context(), ArticleCtxKey, article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Order middleware is used to load an Order object from
// the URL parameters passed through as the request. In case
// the Order could not be found, we stop here and return a 404.
func Order(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var order *types.Order

		if id := chi.URLParam(r, "id"); id != "" {
			intID, err := strconv.Atoi(id)
			if err != nil {
				_ = render.Render(w, r, types.ErrInvalidRequest(err))
				return
			}
			order = DBClient.GetOrderByID(intID)
		} else {
			_ = render.Render(w, r, types.ErrNotFound())
			return
		}
		if order == nil {
			_ = render.Render(w, r, types.ErrNotFound())
			return
		}

		ctx := context.WithValue(r.Context(), OrderCtxKey, order)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
