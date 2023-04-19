// server.go tempat dimana kita akan mengimplementasikan http api server
package api

import (
	db "backendmaster/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves all http request for our banking service
type Server struct {
	store db.Store

	// router mengirimkan setiap api request ke api handler yang tepat untuk diproses
	router *gin.Engine
}

// NewServer membuat instansi baru dari Server
// and setup all api route untuk semua service di Server tsb
func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()
	// mendaftarkan custom validator ke gin
	// get current validator engine yang gin gunakan, konversi outputnya menjadi (*validator.validate)
	if v,ok := binding.Validator.Engine().(*validator.Validate); ok {
		//v.RegisterValidation("name validation tag", function)
		v.RegisterValidation("currency", validCurrency) // untuk mendaftarkan custom validator buatan kita
	} 

	// menambahkan route pada router
	// jika kita mengirimkan multiple function pada handlernya
	//  maka function yang urutannya terakhir akan menjadi handler yang asli
	// function sebelumnya hanya menjadi middle ware
	// post method untuk membuat account baru
	router.POST("/accounts", server.createAccount)

	router.GET("/accounts/:id", server.getAccount)

	// tidak perlu uri, kita akan mendapatkan req data dari query params
	router.GET("/accounts", server.listAccount)

	router.PUT("/accounts", server.updateAccount)

	router.DELETE("/accounts", server.deleteAccount)

	router.POST("/transfers", server.CreateTransfer)

	router.POST("/users", server.CreateUser)


	server.router = router
	return server
}

// start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
	// LATER kita bisa mengimplementasikan gracefully shadow logic di function ini
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
	// LATER bisa implementasikan untuk mengembalikan sesuai tipe errornya
}
