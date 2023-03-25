package main

import (
	"backendmaster/api"
	db "backendmaster/db/sqlc"
	"backendmaster/utils/conf"
	"database/sql"
	"log"

	_ "github.com/lib/pq" // tanpa ini kode tidak akan bisa berkomunikasi dengan database
)

// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://root:21204444@localhost:5432/bank_mandiri?sslmode=disable" // copy saja dari migrate command
// 	serverAddress = "0.0.0.0:8080"
// )



func main() {
	// mengambil config yang sudah diberikan oleh viper
	config,err := conf.LoadConfig(".") // membaca file config dilokasi yang sama , lokasi cukup sampai pada foler yang nampung app.env saja, app.env tidak dituliska
	if err != nil {
		log.Fatal("tidak bisa membaca configuration : ", err)
	}

	// connect ke database
	// establish connection to  database
	conn, err := sql.Open(config.DBDriver, config.DBSource) // create connection to db
	if err != nil {
		log.Fatal("tidak bisa konek ke database : ", err)
	}

	// create store
	store := db.NewStore(conn)

	server := api.NewServer(store)


	// start server
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("tidak bisa memulai server : ", err)
	}

}
