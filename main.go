package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/techschool/simplebank/api"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}
	conn, err := sql.Open(config.DbDriver, config.DbSource)

	if err != nil {
		log.Fatal("could not connect to db", err)
	}
	//kết nối máy chủ
	store := db.NewStore(conn)
	//tạo máy chủ mới
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Error creating server: ", err)
	}
	err = server.Start(config.ServerDriver)
	if err != nil {
		log.Fatal("could not start server", err)
	}

}
