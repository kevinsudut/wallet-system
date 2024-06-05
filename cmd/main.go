package main

import (
	"errors"
	"net/http"

	"github.com/kevinsudut/wallet-system/app"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
)

func main() {
	log.Init()

	err := app.Init()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln("app.Init", err)
	}
}
