package config

import (
	"github.com/sirupsen/logrus"
	gcfg "gopkg.in/gcfg.v1"
)

// Storage contains needed access keys and enpoints for accessing
// the video and thumbnail storage service.
type StorageConfiguration struct {
	// Endpoint is the URL that will be used to access object storage.
	Endpoint string
	// AccessKeyID is the access key used for object storage.
	AccessKeyID string
	// SecretAccessKey is the secret key used for object storage.
	SecretAccessKey string
	// UseSSL controls whether or not SSL is used to connect to object storage.
	UseSSL bool
	// Database is the string of options used to connect to a relational database.
	Database string
	// DatabaseLog controls the log level of database operations.
	DatabaseLog bool
}

// RenderingConfiguration specifies the directory template files
// will be located in.
type RenderingConfiguration struct {
	// TemplatesDir is the directory that template files can be found int.
	// Templates must end with the .tmpl file extension.
	TemplatesDir string
}

// SessionConfiguration defines the key that is used to encrypt all
// user sessions.
type SessionConfiguration struct {
	// Key used for session encryption.
	Key string
}

// Configuration encapsulates configuration of all submodules within the
// application.
type Configuration struct {
	// Storage contains all configuration for storage services.
	Storage StorageConfiguration
	// Rendering contains all configuration for rendering services.
	Rendering RenderingConfiguration
	// Session contrains all configuration for session management.
	Session SessionConfiguration
}

var conf Configuration

func init() {
	logrus.Info("Parsing configuration...")
	err := gcfg.ReadFileInto(&conf, "config.ini")
	if err != nil {
		logrus.Fatal("Could not read configuration : ", err)
	}
}

// GetStorageConfiguration returns the storage configuration specified in
// the config.ini file
func GetStorageConfiguration() StorageConfiguration {
	return conf.Storage
}

// GetRenderingConfiguration returns the Rednering storage configuration
// specified in the config.ini file
func GetRenderingConfiguration() RenderingConfiguration {
	return conf.Rendering
}

// GetSessionConfiguration returns the session configuration specified in
// the config.ini file
func GetSessionConfiguration() SessionConfiguration {
	return conf.Session
}
