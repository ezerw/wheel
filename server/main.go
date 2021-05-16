package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/ezerw/wheel/db"
	"github.com/ezerw/wheel/handler"
	"github.com/ezerw/wheel/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open("mysql", config.DSN)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server, err := handler.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
