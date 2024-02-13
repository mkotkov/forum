package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"forum/internal/application"
	"forum/internal/repository"
)

func main() {
	ctx := context.Background()

	db, err := repository.InitDBConn(ctx)
	if err != nil {
		log.Fatalf("%v failed to init DB connection", err)
	}
	defer db.Close()

	a := application.NewApp(ctx, db)
	router := &RouterAdapter{a}

	http.Handle("/", router)

	// static file handler
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: nil}

	go func() {
		fmt.Println("It is alive! Try http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	fmt.Println("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	fmt.Println("Server stopped gracefully")
}

// RouterAdapter adapter for converting *application.App to http.Handler interface
type RouterAdapter struct {
	app *application.App
}

// ServeHTTP implements the http.Handler interface method for RouterAdapter
func (ra *RouterAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ra.app.Routes(w, r)
}
