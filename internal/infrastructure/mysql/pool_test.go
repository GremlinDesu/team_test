package pool_test

import (
	"log/slog"
	"team-test/config"
	pool "team-test/internal/infrastructure/mysql"
	"testing"
)

var NewCfg = config.Config{MySQLdsn: "user:password@tcp(localhost:3306)/teams_db"}

func TestSimple(t *testing.T) {
	_, err := pool.NewPool(NewCfg)
	if err != nil {
		slog.Debug(err.Error())
		t.Error(err.Error())
	}
}
