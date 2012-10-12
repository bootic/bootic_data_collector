package db

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	//"twitter1/vendor/redigo/redis"
)

var pg *sql.DB

func prepareOrPanic(query string) (stmt *sql.Stmt) {

  var err error

  if stmt, err = pg.Prepare(query); err != nil {
    panic(err)
  } 

  return stmt
}

func Init() {

	if c, err := sql.Open("postgres", "dbname=datagram sslmode=disable"); err == nil {
		pg = c
	} else {
		panic(err)
	}

  prepareEventStatements()
}
