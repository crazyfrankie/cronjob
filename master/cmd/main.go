package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crazyfrankie/cronjob/master/conf"
	"github.com/crazyfrankie/cronjob/master/ioc"
)

func main() {
	engine := ioc.InitServer()

	srv := &http.Server{
		Addr:    conf.GetConf().Server.Port,
		Handler: engine,
	}

	log.Printf("Server is running at http://localhost:%s", conf.GetConf().Server.Port)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("failed start server")
	}

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown main server: %v", err)
	}
	log.Println("Server exited gracefully")
}
