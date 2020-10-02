package db

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"

	"github.com/jonnylangefeld/go-api/pkg/types"
)

var (
	testClient  = &Client{}
	testArticle = types.Article{
		Name:  "Skittles",
		Price: 1.99,
	}
)

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	r, err := pool.Run("postgres", "13-alpine", []string{"POSTGRES_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		if err := testClient.Connect(fmt.Sprintf("host=localhost port=%s user=postgres dbname=postgres password=secret sslmode=disable", r.GetPort("5432/tcp"))); err != nil {
			return err
		}
		return testClient.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(r); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestClient_Articles(t *testing.T) {
	testClient.Client.DropTable(&types.Article{})
	testClient.autoMigrate()
	first := testArticle
	err := testClient.SetArticle(&first)
	assert.NoError(t, err)
	assert.Equal(t, 1, first.ID)

	second := testArticle
	err = testClient.SetArticle(&second)
	assert.NoError(t, err)
	assert.Equal(t, 2, second.ID)

	update := first
	update.Price = 2.99
	err = testClient.SetArticle(&update)
	assert.NoError(t, err)

	got := testClient.GetArticleByID(1)
	assert.Equal(t, testArticle.Name, got.Name, "")
	assert.Equal(t, 2.99, got.Price, "")

	got = testClient.GetArticleByID(2)
	assert.Equal(t, testArticle.Name, got.Name, "")
	assert.Equal(t, 1.99, got.Price, "")
}

func TestClient_PaginateArticles(t *testing.T) {
	testClient.Client.DropTable(&types.Article{})
	testClient.autoMigrate()
	for i := 0; i < pageSize+2; i++ {
		article := testArticle
		_ = testClient.SetArticle(&article)
	}
	got := testClient.GetArticles(0)
	assert.Equal(t, 10, len(got.Items))
	assert.Equal(t, 11, got.NextPageID)

	got = testClient.GetArticles(11)
	assert.Equal(t, 2, len(got.Items))
	assert.Equal(t, 0, got.NextPageID)
}
