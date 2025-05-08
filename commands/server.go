package commands

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/osesantos/resulto"
)

func Serve() resulto.ResultAny {
	fmt.Println("Server is running...")

	router := gin.Default()

	// /ping
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// /
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	router.Run()

	return resulto.SuccessAny()
}
