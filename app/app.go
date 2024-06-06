package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/kevinsudut/wallet-system/app/handler"
	"github.com/kevinsudut/wallet-system/pkg/lib/database"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
	"github.com/kevinsudut/wallet-system/pkg/lib/redis"
	"github.com/kevinsudut/wallet-system/pkg/lib/token"
)

func Init() error {
	db, err := database.Init()
	if err != nil {
		return err
	}

	redis, err := redis.Init()
	if err != nil {
		return err
	}

	token, err := token.Init()
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr: ":8000",
		Handler: http.TimeoutHandler(
			handler.Init(token, db, redis).RegisterHandlers(mux.NewRouter()),
			1*time.Second, // 1 second as timeout
			"",
		),
		ReadTimeout:  1 * time.Second, // 1 second as timeout
		WriteTimeout: 1 * time.Second, // 1 second as timeout
	}

	// Terminated gracefully using SIGTERM
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		if err := server.Shutdown(context.TODO()); err != nil {
			log.Errorln("Error shutting down server", err)
		}
	}()

	return server.ListenAndServe()
}
