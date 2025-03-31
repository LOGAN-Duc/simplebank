package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/techschool/simplebank/Utill"
	"github.com/techschool/simplebank/api"
	db "github.com/techschool/simplebank/db/sqlc"
	"log"
)

func main() {
	config, err := Utill.LoadConfig(".")
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
	server := api.NewServer(store)

	err = server.Start(config.ServerDriver)
	if err != nil {
		log.Fatal("could not start server", err)
	}

}
