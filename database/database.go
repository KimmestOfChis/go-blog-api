package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

const (
	host    = "localhost"
	port    = 5432
	user   = "kimcaraway"
	password = "SmallKimjay_2008"
	dbname  = "go_blog_api"
)

func Connect() (*sql.DB, error) {
	connStr := "user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	return sql.Open("postgres", connStr)
}

