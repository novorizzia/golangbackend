package api

import (
	db "backendmaster/db/sqlc"
	"backendmaster/utils/password"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	// alphanum = membuat user tidak dapat menginputkan karakter untuk username selain alphabet dan number
	Username string `json:"username" binding:"required,alphanum"`
	// panjang password sebaiknya tidak terlalu pendek
	Password string `json:"password" binding:"required,min=4"`
	FullName string `json:"full_name" binding:"required"`
	// email = menvalidasi agar data yang diinputkan user berformat email yang benar
	Email string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	Fullname          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {

	return userResponse{
		Username:          user.Username,
		Fullname:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) CreateUser(ctx *gin.Context) {
	// context objek menyediakan method untuk menulis input parameter dan mengembalikan responses

	var req createUserRequest

	// mengikat parameter body ke variable req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// jika terjadi error maka kemungkinan user menginputkan data yang invalid
		// jadi kita harus mengirimkan bad request respon pada client
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := password.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// print untuk melihat error code name
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		} // konversi error ke tipe pq.error

		// jika error tidak nill pasti terdapat isu internal
		// jadi kita harus mengirimkan status bahwa server sedang error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := newUserResponse(user)

	ctx.JSON(http.StatusOK, res)

}

type loginUserRequest struct {
	// alphanum = membuat user tidak dapat menginputkan karakter untuk username selain alphabet dan number
	Username string `json:"username" binding:"required,alphanum"`
	// panjang password sebaiknya tidak terlalu pendek
	Password string `json:"password" binding:"required,min=4"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	User        userResponse
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = password.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	res := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, res)

}
