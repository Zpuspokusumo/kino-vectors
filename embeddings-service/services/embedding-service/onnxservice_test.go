package embeddingservice

import (
	"fmt"
	"kino-vectors/env"
	"kino-vectors/repository"
	"log"
	"os"
	"testing"
)

func setupservice(ENV env.ENV) (*EmbeddingServiceONNX, error) {
	client, err := repository.NewClient(ENV.QdrantAPIKEY)
	if err != nil {
		return nil, err
	}

	repo := repository.QdrantRepository{Client: client}

	service, err := MakeServiceONNX(repo)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func TestEmbedding(t *testing.T) {
	ENV := env.Setup()

	service, err := setupservice(ENV)
	if err != nil {
		fmt.Println(err)
		log.Fatal("cant setup service, shutting down")
	}

	data, err := service.GenerateEmbeddings("text here")
	if err != nil {
		fmt.Println(err)
		log.Fatal("cant run, shutting down")
	}

	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	// The defer statement ensures the file is closed at the end of the main function.
	defer file.Close()

	for _, number := range data {
		_, err := fmt.Fprintf(file, "%.5f\n", number)
		if err != nil {
			log.Fatalf("failed to write to file: %v", err)
		}
	}
}
