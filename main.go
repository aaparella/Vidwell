package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aaparella/vidwell/storage"
	"github.com/gorilla/mux"
)

// Teardown is called when the application terminates
func Teardown(signal os.Signal) {
	log.Println("Received", signal, "signal, shutting down")
	log.Println("Database teardown...")
	if err := storage.Teardown(); err != nil {
		log.Println("Storage teardown error : ", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	ch := make(chan os.Signal)
	go func() {
		Teardown(<-ch)
	}()
	signal.Notify(ch, syscall.SIGINT, syscall.SIGKILL)

	router := mux.NewRouter()
	RegisterRoutes(router)
	http.ListenAndServe(":8080", router)
}
