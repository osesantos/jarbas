package commands

import (
	"fmt"

	"github.com/osesantos/resulto"
)

func Serve() resulto.ResultAny {
	fmt.Println("Server is running...")

	return resulto.SuccessAny()
}
