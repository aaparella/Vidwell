package videos

import (
	"log"

	"code.google.com/p/go-uuid/uuid"

	"github.com/aaparella/vidwell/storage"
)

// UploadVideo uploads video content to content storage, nothing else
func UploadVideo(data []byte, contentType string) {
	if err := storage.Upload(data, uuid.New()+".mov", "videos", contentType); err != nil {
		log.Println(err)
	}
}
