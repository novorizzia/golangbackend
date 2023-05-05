// server.go tempat dimana kita akan mengimplementasikan http api server
package api

import (
	db "backendmaster/db/sqlc"
	"backendmaster/token"
	"backendmaster/utils/conf"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves all http request for our banking service
type Server struct {
	config conf.Config
	store  db.Store

	tokenMaker token.Maker

	// router mengirimkan setiap api request ke api handler yang tepat untuk diproses
	router *gin.Engine
}

// NewServer membuat instansi baru dari Server
// and setup all api route untuk semua service di Server tsb
func NewServer(conf conf.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(conf.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("tidak bisa membuat token maker : %w", err)
	}

	server := &Server{
		config:     conf,
		store:      store,
		tokenMaker: tokenMaker,
		router:     &gin.Engine{},
	}

	// mendaftarkan custom validator ke gin
	// get current validator engine yang gin gunakan, konversi outputnya menjadi (*validator.validate)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//v.RegisterValidation("name validation tag", function)
		v.RegisterValidation("currency", validCurrency) // untuk mendaftarkan custom validator buatan kita
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	// menambahkan route pada router
	// jika kita mengirimkan multiple function pada handlernya
	//  maka function yang urutannya terakhir akan menjadi handler yang asli
	// function sebelumnya hanya menjadi middle ware
	// post method untuk membuat account baru

	router.POST("/users", server.CreateUser)
	router.POST("/users/login", server.loginUser)

	// "/" path prefix
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// daripada router. kita gunkan authRoutes untuk memasukan kedalam group route. skrng semua ruoute di group tersebut akan berbagi middleware yang sama
	// smua request yang masuk rute ini akan melewati middleware terlebih dahulu
	authRoutes.POST("/accounts", server.createAccount)

	authRoutes.GET("/accounts/:id", server.getAccount)

	// tidak perlu uri, kita akan mendapatkan req data dari query params
	authRoutes.GET("/accounts", server.listAccount)

	authRoutes.PUT("/accounts", server.updateAccount)

	authRoutes.DELETE("/accounts", server.deleteAccount)

	authRoutes.POST("/transfers", server.CreateTransfer)

	server.router = router
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
