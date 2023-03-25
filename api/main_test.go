package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode) // set gin to test mode
	os.Exit(m.Run()) // to start unit test, mengembalikan pass atau fail
}
