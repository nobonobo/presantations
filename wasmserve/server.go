package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nobonobo/wasmserve/util"
)

func main() {
	util.SetWorkdir(".")
	log.Println("server start: http://localhost:5000")
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	server := &http.Server{
		Addr:    ":5000",
		Handler: http.DefaultServeMux,
	}
	go func() {
		<-quit
		log.Println("Server is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	<-done
	log.Println("Server stopped")
	time.Sleep(time.Second)
}
