package storage

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/aaparella/vidwell/config"
	minio "github.com/minio/minio-go"
)

var client *minio.Client

func InitializeObjectStorage(conf config.StorageConfiguration) error {
	var err error
	client, err = minio.New(conf.Endpoint,
		conf.AccessKeyID,
		conf.SecretAccessKey,
		conf.UseSSL)
	if err != nil {
		return fmt.Errorf("\nCould not connect to storage service: %s", err)
	}

	if err := CreateBuckets("vidwell.videos", "vidwell.thumbnails", "vidwell.avatars"); err != nil {
		return fmt.Errorf("Error creating necessary buckets: %s", err)
	}

	return nil
}

// Upload creates an object with the specified name, in the specified bucket,
// with the passed data and content type. This is the base function that will
// be used for all data uploading, be it thumbnails or avatars or whatever
// else we want uploaded.
func Upload(data []byte, name, bucket, contentType string) error {
	reader := bytes.NewReader(data)

	_, err := client.PutObject(bucket, name, reader, contentType)
	return err
}

func CreateBuckets(names ...string) error {
	buckets, err := client.ListBuckets()
	if err != nil {
		return err
	}

	check := func(name string) {
		if err == nil {
			err = EnsureExists(name, buckets)
		}
	}

	for _, name := range names {
		check(name)
	}

	return err
}

func EnsureExists(name string, buckets []minio.BucketInfo) error {
	if exists := BucketExists(name, buckets); !exists {
		return CreateBucket(name)
	}
	return nil
}

func BucketExists(name string, buckets []minio.BucketInfo) bool {
	for _, bucket := range buckets {
		if bucket.Name == name {
			return true
		}
	}
	return false
}

func CreateBucket(name string) error {
	return client.MakeBucket(name, "")
}

func GetObjectUrl(name, bucket string) *url.URL {
	url, err := client.PresignedGetObject(bucket, name, time.Hour, nil)
	if err != nil {
		log.Println("Could not get URL : ", err.Error())
		return nil
	}

	return url
}

func GetVideoUrl(uuid string) *url.URL {
	return GetObjectUrl(uuid, "vidwell.videos")
}
