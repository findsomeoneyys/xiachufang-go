package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/findsomeoneyys/xiachufang-api/config"
	"github.com/findsomeoneyys/xiachufang-api/router"
	"github.com/gin-gonic/gin"
)

func main() {

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := config.Get()

	gin.SetMode(config.Server.RunMode)
	g := gin.Default()
	router.InitRouter(g)

	endPoint := fmt.Sprintf(":%d", config.Server.HttpPort)

	server := &http.Server{
		Addr:           endPoint,
		Handler:        g,
		ReadTimeout:    config.Server.ReadTimeout * time.Second,
		WriteTimeout:   config.Server.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")

}
