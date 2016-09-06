package views

import (
	"fmt"
	"net/http"
)

// LandingPage serves Vidwell's homepage. Obviously.
func LandingPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to VidWell")
}
