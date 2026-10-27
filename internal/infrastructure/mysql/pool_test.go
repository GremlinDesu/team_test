package mysql_pool_test

import (
	"log/slog"
	"team-test/config"
	mysql_pool "team-test/internal/infrastructure/mysql"
	"testing"
)

var NewCfg = config.Config{MySQLdsn: "user:password@tcp(127.0.0.1:3306)/teams_db?parseTime=true"}

func TestSimple(t *testing.T) {
	db, err := mysql_pool.NewPool(NewCfg)
	if err != nil {
		slog.Debug(err.Error())
		t.Error(err.Error())
	}
	err = db.Ping()
	if err != nil {
		slog.Debug(err.Error())
		t.Error(err.Error())
	}
}
