package api

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jonnylangefeld/go-api/pkg/types"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/jonnylangefeld/go-api/pkg/api/mocks"
)

var (
	testArticle1 = types.Article{
		ID:    1,
		Name:  "Skittles",
		Price: 1.99,
	}
	testArticle2 = types.Article{
		ID:    2,
		Name:  "Jelly Beans",
		Price: 2.99,
	}
)

// TestGetRouter ensures that the router contains all expected routes
func TestGetRouter(t *testing.T) {
	log, _ := zap.NewProduction(zap.WithCaller(false))
	r := GetRouter(log, nil)

	testcases := map[string]struct {
		method string
		path   string
	}{
		"GET /articles": {
			method: http.MethodGet,
			path:   "/articles",
		},
		"GET /articles/{id}": {
			method: http.MethodGet,
			path:   "/articles/id",
		},
		"GET /orders": {
			method: http.MethodGet,
			path:   "/orders",
		},
		"GET /orders/{id}": {
			method: http.MethodGet,
			path:   "/orders/id",
		},
		"GET /swagger": {
			method: http.MethodGet,
			path:   "/swagger",
		},
	}

	for name, test := range testcases {
		t.Run(name, func(t *testing.T) {
			got := r.Match(chi.NewRouteContext(), test.method, test.path)
			assert.Equal(t, true, got, fmt.Sprintf("not found: %s '%s'", test.method, test.path))
		})
	}
}

//go:generate $GOPATH/bin/mockgen -destination=./mocks/db.go -package=mocks github.com/jonnylangefeld/go-api/pkg/db ClientInterface

func getDBClientMock(t *testing.T) *mocks.MockClientInterface {
	ctrl := gomock.NewController(t)
	dbClient := mocks.NewMockClientInterface(ctrl)

	dbClient.EXPECT().GetArticles(gomock.Eq(0)).Return(&types.ArticleList{
		Items: []*types.Article{
			&testArticle1,
			&testArticle2,
		},
	})

	dbClient.EXPECT().GetArticles(gomock.Eq(1)).Return(&types.ArticleList{
		Items: []*types.Article{
			&testArticle2,
		},
	})

	dbClient.EXPECT().GetArticleByID(gomock.Eq(1)).Return(&testArticle1).AnyTimes()

	dbClient.EXPECT().SetArticle(gomock.Any()).DoAndReturn(func(article *types.Article) error {
		if article.ID == 0 {
			article.ID = 1
		}
		return nil
	}).AnyTimes()

	return dbClient
}

// TestEndpoints ensures the expected results upon requests
func TestEndpoints(t *testing.T) {
	r := GetRouter(nil, getDBClientMock(t))
	ts := httptest.NewServer(r)
	defer ts.Close()

	testcases := map[string]struct {
		method   string
		path     string
		body     string
		header   http.Header
		wantCode int
		wantBody string
	}{
		"GET /articles": {
			method:   http.MethodGet,
			path:     "/articles",
			wantCode: http.StatusOK,
			wantBody: `{"items":[{"id":1,"name":"Skittles","price":1.99},{"id":2,"name":"Jelly Beans","price":2.99}]}`,
		},
		"GET /articles?page_id=1": {
			method:   http.MethodGet,
			path:     "/articles?page_id=1",
			wantCode: http.StatusOK,
			wantBody: `{"items":[{"id":2,"name":"Jelly Beans","price":2.99}]}`,
		},
		"GET /articles/{id}": {
			method:   http.MethodGet,
			path:     "/articles/1",
			wantCode: http.StatusOK,
			wantBody: `{"id":1,"name":"Skittles","price":1.99}`,
		},
		"PUT /articles": {
			method: http.MethodPut,
			path:   "/articles",
			body:   `{"name":"Skittles","price":1.99}`,
			header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantCode: http.StatusOK,
			wantBody: `{"id":1,"name":"Skittles","price":1.99}`,
		},
		"Page Not Found": {
			method:   http.MethodGet,
			path:     "/blah",
			wantCode: http.StatusNotFound,
			wantBody: "404 page not found",
		},
	}

	for name, test := range testcases {
		t.Run(name, func(t *testing.T) {
			body := bytes.NewReader([]byte(test.body))
			gotResponse, gotBody := testRequest(t, ts, test.method, test.path, body, test.header)
			assert.Equal(t, test.wantCode, gotResponse.StatusCode)
			if test.wantBody != "" {
				assert.Equal(t, test.wantBody, gotBody, "body did not match")
			}
		})
	}
}

// testRequest is a helper function to exectute the http request against the server
func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader, header http.Header) (*http.Response, string) {
	t.Helper()
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	req.Header = header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	respBody = bytes.TrimSpace(respBody)

	return resp, string(respBody)
}
