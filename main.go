/*
In this package we set the specs of the service and define the handlers.
For simplicity I only defined a GET
*/

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"redis/handlers"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "nominatim-api", log.LstdFlags)

	nh := handlers.NewApi(l)

	sm := mux.NewRouter()

	GetRouter := sm.Methods(http.MethodGet).Subrouter()
	GetRouter.HandleFunc("/api", nh.RetrieveNominatim)

	s := http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Println("Starting server on port ", os.Getenv("PORT"))
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	l.Println("Received terminate, graceful shutdown", sig)
	tc, err := context.WithTimeout(context.Background(), 30*time.Second)

	if err != nil {
		l.Fatal(err)
	}

	s.Shutdown(tc)

}
