package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aaparella/vidwell/models"
	"github.com/aaparella/vidwell/storage"
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

	go videos.StoreVideo(data,
		handler.Header.Get("Content-Type"),
		r.FormValue("title"),
		"aaparella")

	fmt.Fprintf(w, "Thank you for the video!")
}

const VIDEOS_HTML = `
	{{ range . }}
		<div class="video"> 
			<h1>{{ .Title }} {{ .Uuid }}</h1>
		</div>
	{{ end }}
`

func Videos(w http.ResponseWriter, r *http.Request) {
	var videos []models.Video
	if err := storage.DB.Find(&videos).Error; err != nil {
		fmt.Fprintf(w, "Could not find videos : %s", err.Error())
		return
	}

	tmpl, _ := template.New("test").Parse(VIDEOS_HTML)
	tmpl.Execute(w, videos)
}

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

	http.HandleFunc("/upload", Upload)
	http.HandleFunc("/videos", Videos)
	http.ListenAndServe(":8080", nil)
}
