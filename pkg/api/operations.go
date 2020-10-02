package api

import (
	"net/http"

	"github.com/go-chi/render"

	m "github.com/jonnylangefeld/go-api/pkg/middelware"
	"github.com/jonnylangefeld/go-api/pkg/types"
)

// GetArticle renders the article from the context
// @Summary Get article by id
// @Description GetArticle returns a single article by id
// @Tags Articles
// @Produce json
// @Param id path string true "article id"
// @Router /articles/{id} [get]
// @Success 200 {object} types.Article
// @Failure 400 {object} types.ErrResponse
// @Failure 404 {object} types.ErrResponse
func GetArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value(m.ArticleCtxKey).(*types.Article)

	if err := render.Render(w, r, article); err != nil {
		_ = render.Render(w, r, types.ErrRender(err))
		return
	}
}

// PutArticle writes an article to the database
// @Summary Add an article to the database
// @Description PutArticle writes an article to the database
// @Description To write a new article, leave the id empty. To update an existing one, use the id of the article to be updated
// @Tags Articles
// @Produce json
// @Router /articles [put]
// @Success 200 {object} types.Article
// @Failure 400 {object} types.ErrResponse
// @Failure 404 {object} types.ErrResponse
func PutArticle(w http.ResponseWriter, r *http.Request) {
	article := &types.Article{}
	if err := render.Bind(r, article); err != nil {
		_ = render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}

	if err := DBClient.SetArticle(article); err != nil {
		_ = render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}

	if err := render.Render(w, r, article); err != nil {
		_ = render.Render(w, r, types.ErrRender(err))
		return
	}
}

// ListArticles returns all articles in the database
// @Summary List all articles
// @Description Get all articles stored in the database
// @Tags Articles
// @Produce json
// @Param page_id query string false "id of the page to be retrieved"
// @Router /articles [get]
// @Success 200 {object} types.ArticleList
// @Failure 400 {object} types.ErrResponse
// @Failure 404 {object} types.ErrResponse
func ListArticles(w http.ResponseWriter, r *http.Request) {
	pageID := r.Context().Value(m.PageIDKey)
	if err := render.Render(w, r, DBClient.GetArticles(pageID.(int))); err != nil {
		_ = render.Render(w, r, types.ErrRender(err))
		return
	}
}

// GetOrder renders the order from the context
// @Summary Get order by id
// @Description GetOrder returns a single order by id
// @Tags Orders
// @Produce json
// @Param id path string true "order id"
// @Router /orders/{id} [get]
// @Success 200 {object} types.Order
// @Failure 400 {object} types.ErrResponse
// @Failure 404 {object} types.ErrResponse
func GetOrder(w http.ResponseWriter, r *http.Request) {
	order := r.Context().Value(m.OrderCtxKey).(*types.Order)

	if err := render.Render(w, r, order); err != nil {
		_ = render.Render(w, r, types.ErrRender(err))
		return
	}
}

// PutOrder writes an order to the database
// @Summary Add an order to the database
// @Description PutOrder writes an order to the database
// @Description To write a new order, leave the id empty. To update an existing one, use the id of the order to be updated
// @Tags Orders
// @Produce json
// @Router /orders [put]
// @Success 200 {object} types.Order
// @Failure 400 {object} types.ErrResponse
// @Failure 404 {object} types.ErrResponse
func PutOrder(w http.ResponseWriter, r *http.Request) {
	order := &types.Order{}
	if err := render.Bind(r, order); err != nil {
		_ = render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}

	if err := DBClient.SetOrder(order); err != nil {
		_ = render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}

	if err := render.Render(w, r, order); err != nil {
		_ = render.Render(w, r, types.ErrRender(err))
		return
	}
}

// ListOrders returns all orders in the database
// @Summary List all orders
// @Description Get all orders stored in the database
// @Tags Orders
// @Produce json
// @Param page_id query string false "id of the page to be retrieved"
// @Router /orders [get]
// @Success 200 {object} types.OrderList
// @Failure 400 {object} types.ErrResponse
// @Failure 404 {object} types.ErrResponse
func ListOrders(w http.ResponseWriter, r *http.Request) {
	pageID := r.Context().Value(m.PageIDKey)
	if err := render.Render(w, r, DBClient.GetOrders(pageID.(int))); err != nil {
		_ = render.Render(w, r, types.ErrRender(err))
		return
	}
}
