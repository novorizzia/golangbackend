package db

import (
	"backendmaster/utils/conf"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

// untuk sementara kita gunakan konstanta, dalam real case kita akan menarik data dari environment variable
// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://root:21204444@localhost:5432/bank_mandiri?sslmode=disable" // copy saja dari migrate command
// )

// agar dapat digunakan di file test lainnya di package yang sama
var testDB *sql.DB

// special function whic entry point to run all unit test inside one spesific golang package
func TestMain(m *testing.M) {
	config, err := conf.LoadConfig("./../../") // lokasi cukup sampai pada foler yang nampung app.env saja, app.env tidak dituliskan
	if err != nil {
		log.Fatal("tidak bisa membaca configuration : ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource) // create connection to db

	if err != nil {
		log.Fatal("tidak bisa konek ke database : ", err)
	}

	testQueries = New(testDB) // function new dari file  yang digen sqlc

	os.Exit(m.Run()) // to start unit test, mengembalikan pass atau fail
}
