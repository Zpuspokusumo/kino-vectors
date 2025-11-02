package main

import (
	"context"
	"fmt"
	"kino-vectors/data"
	"kino-vectors/env"
	"kino-vectors/repository"
	embeddingservice "kino-vectors/services/embedding-service"
	"log"
	"os"
)

const (
	maxMessageSize = 20 * 1024 * 1024 // 20MB
)

func setupservice(ENV env.ENV) (*embeddingservice.EmbeddingServiceONNX, error) {
	client, err := repository.NewClient(ENV.QdrantAPIKEY)
	if err != nil {
		return nil, err
	}

	repo := repository.QdrantRepository{Client: client}

	service, err := embeddingservice.MakeServiceONNX(repo)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func main() {

	ENV := env.Setup()

	service, err := setupservice(ENV)
	if err != nil {
		fmt.Println(err)
		log.Fatal("cant setup service, shutting down")
	}
	//text := "I am a witch and i like casting spells"
	text := "Hello, world!"
	testonnxembeddingtofile(service, text)

	//testqdrantinsert(context.Background(), service, ENV.QdrantMovieCollection)
}

func testonnxembeddingtofile(service *embeddingservice.EmbeddingServiceONNX, text string) {

	data, err := service.GenerateEmbeddings(text)
	if err != nil {
		fmt.Println(err)
		log.Fatal("cant run, shutting down")
	}

	file, err := os.OpenFile("data2.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func testqdrantinsert(ctx context.Context, service *embeddingservice.EmbeddingServiceONNX, coll string) {
	movdata := data.GetdataShort()

	data, err := service.ProcessMovieData(ctx, movdata, coll)
	if err != nil {
		fmt.Println(err)
		log.Fatal("cant run, shutting down")
	}
	for i := range data {
		fmt.Printf("id %v score %v\n payload %v\n", data[i].Id, data[i].Score, data[i].Payload)
		for key, value := range data[i].Payload {
			fmt.Printf("Key: %s, Value: %v\n", key, value.String())
		}
		fmt.Printf("\n\n")
	}
}
