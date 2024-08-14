package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Engine   string
	User     string
	Password string
	Protocol string
	Host     string
	Port     string
	Database string
}

func New(c *Config) *sql.DB {
	db, err := sql.Open(
		c.Engine,
		fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true", c.User, c.Password, c.Protocol, c.Host, c.Port, c.Database),
	)

	if err != nil {
		log.Fatal(err)
	}
	return db
}
