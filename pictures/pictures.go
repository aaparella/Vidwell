package pictures

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/aaparella/vidwell/storage"
)

// UploadPicture is the primitive picture uploading function used to upload
// pictures to any desired bucket. Called by all other "Upload..." functions.
func UploadPicture(data []byte, bucket string, contentType string) error {
	return storage.Upload(data, uuid.New(), bucket, contentType)
}

// UploadThumbnail uploads a copy of the passed data to the thumbnails bucket.
func UploadThumbnail(data []byte, contentType string) error {
	return UploadPicture(data, "thumbnails", contentType)
}

// UploadAvatar uploads a copy of the passed data to the avatars bucket.
func UploadAvatar(data []byte, contentType string) error {
	return UploadPicture(data, "avatars", contentType)
}
