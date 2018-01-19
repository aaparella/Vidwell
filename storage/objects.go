package storage

import (
	"bytes"
	"fmt"
	"net/url"
	"time"

	"github.com/aaparella/vidwell/config"
	minio "github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
)

var client *minio.Client

// InitializeObjectStorage creates the requisite buckets for storage if they do
// not already exist, and ensures that the configuration provided is valid and
// allows us to connect to S3 or whatever the service may be.
func InitializeObjectStorage(conf config.StorageConfiguration) error {
	var err error
	client, err = minio.New(conf.Endpoint,
		conf.AccessKeyID,
		conf.SecretAccessKey,
		conf.UseSSL)
	if err != nil {
		return fmt.Errorf("\nCould not connect to storage service: %s", err)
	}

	if err := CreateBuckets("vidwell.videos",
		"vidwell.thumbnails",
		"vidwell.avatars"); err != nil {
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

// CreateBuckets creates buckets with the names passed in to the function.
func CreateBuckets(names ...string) error {
	buckets, err := client.ListBuckets()
	if err != nil {
		return err
	}

	for _, name := range names {
		if err == nil {
			err = EnsureExists(name, buckets)
		}
	}

	return err
}

// EnsureExists checks if the named bucket exists, and create it if it doesn't.
func EnsureExists(name string, buckets []minio.BucketInfo) error {
	if exists := BucketExists(name, buckets); !exists {
		return CreateBucket(name)
	}
	return nil
}

// BucketExists checks if a bucket with the specified name is exists.
func BucketExists(name string, buckets []minio.BucketInfo) bool {
	for _, bucket := range buckets {
		if bucket.Name == name {
			return true
		}
	}
	return false
}

// CreateBucket creates a bucket in the object storage service that
// we have already connected to. If the client is not yet connected,
// this will fail.
func CreateBucket(name string) error {
	return client.MakeBucket(name, "")
}

// GetObjectUrl gets a public facing URL that is valid for one hour, and allows
// us to stream the contents of that bucket to a user. Needs to be reworked.
func GetObjectUrl(name, bucket string) *url.URL {
	url, err := client.PresignedGetObject(bucket, name, time.Hour, nil)
	if err != nil {
		logrus.Println("Could not get URL : ", err.Error())
		return nil
	}

	return url
}
