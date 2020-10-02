package api

import (
	"net/http"

	"github.com/jonnylangefeld/go-api/pkg/db"
	"go.uber.org/zap"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	m "github.com/jonnylangefeld/go-api/pkg/middelware"
)

var DBClient db.ClientInterface

func SetDBClient(c db.ClientInterface) {
	DBClient = c
	m.SetDBClient(DBClient)
}

// GetRouter configures a chi router and starts the http server
// @title My API
// @description This API is a sample go-api.
// @description It also does this.
// @contact.name Jonny Langefeld
// @contact.email jonny.langefeld@gmail.com
// @host example.com
// @BasePath /
func GetRouter(log *zap.Logger, dbClient db.ClientInterface) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	SetDBClient(dbClient)
	if log != nil {
		r.Use(m.SetLogger(log))
	}
	buildTree(r)

	return r
}

func buildTree(r *chi.Mux) {
	r.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
	})
	r.Get("/swagger*", httpSwagger.Handler())

	r.Route("/articles", func(r chi.Router) {
		r.With(m.Pagination).Get("/", ListArticles)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(m.Article)
			r.Get("/", GetArticle)
		})

		r.Put("/", PutArticle)
	})

	r.Route("/orders", func(r chi.Router) {
		r.With(m.Pagination).Get("/", ListOrders)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(m.Order)
			r.Get("/", GetOrder)
		})

		r.Put("/", PutOrder)
	})
}
