package types

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

// Article is one instance of an article
type Article struct {
	// The unique id of this item
	ID int `gorm:"type:SERIAL;PRIMARY_KEY" json:"id" example:"1"`
	// The name of this item
	Name string `gorm:"type:varchar;NOT NULL" json:"name" example:"Skittles"`
	// The price of this item
	Price float64 `gorm:"type:decimal;NOT NULL" json:"price" example:"1.99"`
} // @name Article

// Render implements the github.com/go-chi/render.Renderer interface
func (a *Article) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Bind implements the the github.com/go-chi/render.Binder interface
func (a *Article) Bind(r *http.Request) error {
	return nil
}

// ArticleList contains a list of articles
type ArticleList struct {
	// A list of articles
	Items []*Article `json:"items"`
	// The id to query the next page
	NextPageID int `json:"next_page_id,omitempty" example:"10"`
} // @name ArticleList

// Render implements the github.com/go-chi/render.Renderer interface
func (a *ArticleList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Order is one instance of an order
type Order struct {
	// The unique id of this order
	ID int `gorm:"type:SERIAL;PRIMARY_KEY" json:"id" example:"1"`
	// DateTime is the date and time of this order
	DateTime time.Time `gorm:"timestamp" json:"lastUpdated,omitempty" example:"0001-01-01 00:00:00+00"`
} // @name Order

// Render implements the github.com/go-chi/render.Renderer interface
func (o *Order) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Bind implements the the github.com/go-chi/render.Binder interface
func (o *Order) Bind(r *http.Request) error {
	return nil
}

// OrderList contains a list of orders
type OrderList struct {
	// A list of orders
	Items []*Order `json:"items"`
	// The id to query the next page
	NextPageID int `json:"next_page_id,omitempty" example:"10"`
} // @name OrderList

// Render implements the github.com/go-chi/render.Renderer interface
func (o *OrderList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status" example:"Resource not found."`                                         // user-level status message
	AppCode    int64  `json:"code,omitempty" example:"404"`                                                 // application-specific error code
	ErrorText  string `json:"error,omitempty" example:"The requested resource was not found on the server"` // application-level error message, for debugging
} // @name ErrorResponse

// Render implements the github.com/go-chi/render.Renderer interface for ErrResponse
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest returns a structured http response for invalid requests
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// ErrRender returns a structured http response in case of rendering errors
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

// ErrNotFound returns a structured http response if a resource couln't be found
func ErrNotFound() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "Resource not found.",
	}
}
