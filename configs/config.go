package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Database struct {
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

type environments struct {
	Database    Database
	RestPort    string
	GrpcPort    string
	GraphQLPort string
}

func LoadEnvironment() (environments, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Nota: Arquivo .env não encontrado. Usando variáveis de ambiente do SO.")
	}

	return environments{
		Database: Database{
			DBDriver:   envOr("DB_DRIVER", "mysql"),
			DBHost:     envOr("DB_HOST", "localhost"),
			DBPort:     envOr("DB_PORT", "3306"),
			DBUser:     envOr("DB_USER", "root"),
			DBPassword: envOr("DB_PASSWORD", "root"),
			DBName:     envOr("DB_NAME", "orders"),
		},
		RestPort:    envOr("REST_PORT", "3000"),
		GrpcPort:    envOr("GRPC_PORT", "3002"),
		GraphQLPort: envOr("GRAPHQL_PORT", "3001"),
	}, nil
}

func envOr(key, fallback string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}

	return fallback
}
