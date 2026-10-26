package pool

import (
	"database/sql"
	"team-test/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewPool(Config config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", Config.MySQLdsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
