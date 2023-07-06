package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Options struct {
	Host     string
	Port     string
	DBName   string
	Username string
	Password string
}

type Client struct {
	db     *sql.DB
	dbName string
}

func New(ctx context.Context, opt Options) (*Client, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", opt.Username, opt.Password, opt.Host, opt.Port, opt.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &Client{
		db:     db,
		dbName: opt.DBName,
	}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c Client) DB() *sql.DB {
	return c.db
}

func (c Client) RunMigration() error {
	log.Println("Running migration..")
	driver, err := mysql.WithInstance(c.db, &mysql.Config{})
	if err != nil {
		return err
	}

	p, err := filepath.Abs("storage/migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+p,
		c.dbName,
		driver,
	)
	if err != nil {
		return err
	}

	return m.Up()
}
