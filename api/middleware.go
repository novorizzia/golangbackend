package api

import (
	"backendmaster/token"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload" // payload value akan di store di gin context
)

// ini sebenarnya bukan middleware namun hanya higher order function yang akan mengembalikan authentication middleware function
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	// authentication function yang ingin kita implementasikan
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("Authorization header tidak diisi")	
			// function ini membuat kita bisa melakukan abort pada request dan mengirimkan response json pada client
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader) // karna biasanya token akan berbentuk => Bearer dslkfhoehwfuiskfbsdvhseoihffoisdbviehfbfoo. sehingga kita pisahkan dengan acuan white space
		if len(fields) < 2 {
			err := errors.New("format Authorization header tidak valid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("bentuk authorization ini tidak disupport oleh server %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// store payload pada context
		ctx.Set(authorizationPayloadKey, payload)
		// meng forward request ke next handler
		ctx.Next()

	}

}
