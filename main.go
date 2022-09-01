package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/findsomeoneyys/xiachufang-api/config"
	"github.com/findsomeoneyys/xiachufang-api/router"
	"github.com/gin-gonic/gin"
)

func main() {

	config := config.Get()

	gin.SetMode(config.RunMode)
	g := gin.Default()
	router.InitRouter(g)

	endPoint := fmt.Sprintf(":%d", config.HttpPort)

	server := &http.Server{
		Addr:    endPoint,
		Handler: g,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("receive interrupt signal")
		if err := server.Close(); err != nil {
			log.Fatal("Server Close:", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatalf("Server closed unexpect, %s", err.Error())
		}
	}

	log.Println("Server exiting")

}