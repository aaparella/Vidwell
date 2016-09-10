package config

import (
	"log"

	gcfg "gopkg.in/gcfg.v1"
)

// Storage contains needed access keys and enpoints for accessing
// the video and thumbnail storage service.
type StorageConfiguration struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	Database        string
}

// RenderingConfiguration specifies the directory template files
// will be located in.
type RenderingConfiguration struct {
	TemplatesDir string
}

type SessionConfiguration struct {
	Key string
}

// Configuration encapsulates configuration of all submodules within the
// application.
type Configuration struct {
	Storage   StorageConfiguration
	Rendering RenderingConfiguration
	Session   SessionConfiguration
}

var conf Configuration

func init() {
	log.Println("Parsing configuration...")
	err := gcfg.ReadFileInto(&conf, "config.ini")
	if err != nil {
		log.Fatal("Could not read configuration : ", err)
	}
}

func GetStorageConfiguration() StorageConfiguration {
	return conf.Storage
}

func GetRenderingConfiguration() RenderingConfiguration {
	return conf.Rendering
}

func GetSessionConfiguration() SessionConfiguration {
	return conf.Session
}
