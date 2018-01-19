package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aaparella/vidwell/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Teardown is called when the application terminates. Any module that needs to
// perform any logrusic or cleanup on application close needs to be called from
// here.
//
// Any errors will be logrusged but ignored to allow all modules to at least
// attempt to perform teardown steps.
func Teardown(signal os.Signal) {
	logrus.Println("Received", signal, "signal, shutting down")

	teardown := func(name string, fn func() error) {
		logrus.Printf("Tearing down module : " + name + " ...")
		if err := fn(); err != nil {
			fmt.Printf(" ERROR : " + err.Error())
		}
	}

	teardown("Database", storage.TeardownDatabase)
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
