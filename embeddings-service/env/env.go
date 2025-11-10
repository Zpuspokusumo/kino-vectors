package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ENV struct {
	QdrantAPIport         string
	QdrantMovieCollection string
	Embeddingserviceport  string
}

func Setup() ENV {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return ENV{
		QdrantAPIport:         os.Getenv("QDRANT_API_PORT"),
		QdrantMovieCollection: os.Getenv("QDRANT_MOVIE_COLLECTION"),
		Embeddingserviceport:  os.Getenv("EMBEDDINGSERVICEPORT"),
	}
}
