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

	// Обработка всех запросов через RouterAdapter
	http.Handle("/", router)

	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: nil}

	go func() {
		fmt.Println("It is alive! Try http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// Слушаем сигналы для корректного завершения работы сервера
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	fmt.Println("Shutting down the server...")

	// Остановка сервера с таймаутом 5 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	fmt.Println("Server stopped gracefully")
}

// RouterAdapter адаптер для преобразования *application.App в http.Handler интерфейс
type RouterAdapter struct {
	app *application.App
}

// ServeHTTP реализует метод интерфейса http.Handler для RouterAdapter
func (ra *RouterAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ra.app.Routes(w, r)
}
