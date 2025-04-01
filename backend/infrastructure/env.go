package infrastructure

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// Env holds all environment variables in a structured way
type Env struct {
	ServerAddress string
	MongoURI      string
	DBName        string
	SMTPUsername  string
	SMTPPassword  string
	SMTPHost      string
	SMTPPort      string
}

var (
	envInstance *Env
	once        sync.Once
)

// LoadEnv loads environment variables from .env file (existing function)
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// GetEnv gets an environment variable (existing function)
func GetEnv(key string) string {
	return os.Getenv(key)
}

// NewEnv creates and returns a singleton instance of Env
func NewEnv() *Env {
	once.Do(func() {
		LoadEnv() // Ensure environment is loaded

		envInstance = &Env{
			ServerAddress: GetEnv("SERVER_ADDRESS"),
			MongoURI:      GetEnv("MONGO_URI"),
			DBName:        GetEnv("DB_NAME"),
			SMTPUsername:  GetEnv("SMTPUsername"),
			SMTPPassword:  GetEnv("SMTPPassword"),
			SMTPHost:      GetEnv("SMTPHost"),
			SMTPPort:      GetEnv("SMTPPort"),
		}
	})
	return envInstance
}

// GetEnvStruct provides access to the environment variables struct
func GetEnvStruct() *Env {
	if envInstance == nil {
		return NewEnv()
	}
	return envInstance
}