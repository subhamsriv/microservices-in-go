package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/subhamsriv/microservices-in-go/handlers"
)

func main() {
	l := log.New(os.Stdout, "api-name ", log.LstdFlags)
	product := handlers.NewProduct(l)

	sm := mux.NewRouter()
	getRouter := sm.Methods("Get").Subrouter()
	getRouter.HandleFunc("/product", product.GetProducts)

	putRouter := sm.Methods("Put").Subrouter()
	putRouter.HandleFunc("/product/{id:[0-9]+}", product.UpdateProduct)
	putRouter.Use(product.MiddlewareValidateProduct)

	postRouter := sm.Methods("Post").Subrouter()
	postRouter.HandleFunc("/product", product.AddProduct)
	postRouter.Use(product.MiddlewareValidateProduct)

	s := http.Server{
		Addr:         ":8080",
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Println("Starting server at 8080")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)

}
