package mysqlclient

import (
	"database/sql"
	"fmt"

	dbconfig "main/Config"

	_ "github.com/go-sql-driver/mysql"
)

func New(cfg *dbconfig.Config) *sql.DB {
	//connect db here and return db pointer
	var err error
	var db *sql.DB
	s := fmt.Sprintf("%s:%s@(127.0.0.1:3306)/%s", cfg.Config.Username, cfg.Config.Password, cfg.Config.Database)
	db, err = sql.Open(cfg.Config.Host, s)
	if err != nil {
		panic(err.Error())
	}
	db.Ping()
	return db
}
