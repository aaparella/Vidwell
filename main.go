package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aaparella/vidwell/videos"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(2 << 32)
	file, handler, err := r.FormFile("fileupload")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	go videos.UploadVideo(data, handler.Header.Get("Content-Type"))

	fmt.Fprintf(w, "Thank you for the video!")
}

func main() {
	http.HandleFunc("/upload", Upload)
	http.ListenAndServe(":8080", nil)
}
