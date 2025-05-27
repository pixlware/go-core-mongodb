package mongodb

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type MongoConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	URI      string
	URIQuery string
}

var Config MongoConfig = MongoConfig{
	Username: "",
	Password: "",
	Host:     "",
	Port:     "",
	Database: "",
	URIQuery: "",
	URI:      "",
}

func init() {
	env := getEnv("ENV", "default")
	var envFilePath string
	if env == "default" {
		return
	} else {
		envFilePath = ".env." + env
	}

	err := godotenv.Load(envFilePath)
	if err != nil {
		return
	}

	Config.Username = getEnv("MONGO_USERNAME", Config.Username)
	Config.Password = getEnv("MONGO_PASSWORD", Config.Password)
	Config.Host = getEnv("MONGO_HOST", Config.Host)
	Config.Port = getEnv("MONGO_PORT", Config.Port)
	Config.Database = getEnv("MONGO_DATABASE", Config.Database)
	Config.URIQuery = getEnv("MONGO_URI_QUERY", Config.URIQuery)

	var missingConfig []string

	if Config.Host == "" {
		missingConfig = append(missingConfig, "MONGO_HOST")
	}

	if Config.Port == "" {
		missingConfig = append(missingConfig, "MONGO_PORT")
	}

	if Config.Database == "" {
		missingConfig = append(missingConfig, "MONGO_DATABASE")
	}

	if len(missingConfig) > 0 {
		log.Fatal("[MongoDB] Invalid configuration: Missing " + strings.Join(missingConfig, ", "))
		return
	}

	Config.URI = "mongodb://"

	if Config.Username != "" {
		Config.URI += Config.Username
	}

	if Config.Password != "" {
		Config.URI += ":" + Config.Password
	}

	Config.URI += "@" + Config.Host + ":" + Config.Port + "/" + Config.Database

	if Config.URIQuery != "" {
		Config.URI += "?" + Config.URIQuery
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
