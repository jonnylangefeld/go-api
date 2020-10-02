// this package contains all database operations

package db

import (
	"github.com/jinzhu/gorm"
	// postgres blank import for gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jonnylangefeld/go-api/pkg/types"
)

const (
	pageSize = 10
)

// ClientInterface resembles a db interface to interact with an underlying db
type ClientInterface interface {
	Ping() error
	Connect(connectionString string) error
	GetArticleByID(id int) *types.Article
	SetArticle(article *types.Article) error
	GetArticles(pageID int) *types.ArticleList
	GetOrderByID(id int) *types.Order
	SetOrder(order *types.Order) error
	GetOrders(pageID int) *types.OrderList
}

// Client is a custom db client
type Client struct {
	Client *gorm.DB
}

// Ping allows the db to be pinged.
func (c *Client) Ping() error {
	return c.Client.DB().Ping()
}

// Connect establishes a connection to the database and auto migrates the database schema
func (c *Client) Connect(connectionString string) error {
	var err error
	// Create the database connection
	c.Client, err = gorm.Open(
		"postgres",
		connectionString,
	)

	// End the program with an error if it could not connect to the database
	if err != nil {
		return err
	}
	c.Client.LogMode(false)
	c.autoMigrate()
	return nil
}

// autoMigrate creates the default database schema
func (c *Client) autoMigrate() {
	c.Client.AutoMigrate(&types.Article{})
	c.Client.AutoMigrate(&types.Order{})
}

// GetArticleByID queries an article from the database
func (c *Client) GetArticleByID(id int) *types.Article {
	article := &types.Article{}

	c.Client.Where("id = ?", id).First(&article).Scan(article)

	return article
}

// SetArticle writes an article to the database
func (c *Client) SetArticle(article *types.Article) error {
	// Upsert by trying to create and updating on conflict
	if err := c.Client.Create(&article).Error; err != nil {
		return c.Client.Model(&article).Where("id = ?", article.ID).Update(&article).Error
	}
	return nil
}

// GetArticles returns all articles from the database
func (c *Client) GetArticles(pageID int) *types.ArticleList {
	articles := &types.ArticleList{}
	c.Client.Where("id >= ?", pageID).Order("id").Limit(pageSize + 1).Find(&articles.Items)
	if len(articles.Items) == pageSize+1 {
		articles.NextPageID = articles.Items[len(articles.Items)-1].ID
		articles.Items = articles.Items[:pageSize]
	}
	return articles
}

// GetOrderByID queries an order from the database
func (c *Client) GetOrderByID(id int) *types.Order {
	order := &types.Order{}

	c.Client.Where("id = ?", id).First(&order).Scan(order)

	return order
}

// SetOrder writes an order to the database
func (c *Client) SetOrder(order *types.Order) error {
	// Upsert by trying to create and updating on conflict
	if err := c.Client.Create(&order).Error; err != nil {
		return c.Client.Model(&order).Where("id = ?", order.ID).Update(&order).Error
	}
	return nil
}

// GetOrders returns all orders from the database
func (c *Client) GetOrders(pageID int) *types.OrderList {
	orders := &types.OrderList{}
	c.Client.Find(&orders.Items).Where("id >= ?", pageID).Order("id").Limit(pageSize + 1)
	if len(orders.Items) == pageSize+1 {
		orders.NextPageID = orders.Items[len(orders.Items)-1].ID
		orders.Items = orders.Items[:pageSize+1]
	}
	return orders
}
