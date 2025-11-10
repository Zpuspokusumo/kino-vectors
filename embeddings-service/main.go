package main

import (
	"context"
	"fmt"
	"kino-vectors/controller"
	"kino-vectors/data"
	"kino-vectors/env"
	"kino-vectors/repository"
	embeddingservice "kino-vectors/services/embedding-service"
	"log"
	"net"
	"os"

	moviepb "github.com/Zpuspokusumo/kino-vectors/contract/golang/movie-services"
	"google.golang.org/grpc"
)

const (
	maxMessageSize = 20 * 1024 * 1024 // 20MB
)

func setupserver(ENV env.ENV) (*controller.Controller, error) {
	client, err := repository.NewClient(ENV.QdrantAPIport)
	if err != nil {
		return nil, err
	}

	repo := repository.QdrantRepository{Client: client}

	service, err := embeddingservice.MakeServiceONNX(repo)
	if err != nil {
		return nil, err
	}
	c := controller.New(service)
	return c, nil
}

// func setupgrpc() {
// 	a, b := moviepb.RegisterMovieServiceServer()
// }

func main() {

	ENV := env.Setup()

	servercontroller, err := setupserver(ENV)
	if err != nil {
		fmt.Println(err)
		log.Fatal("cant setup service, shutting down")
	}
	//testonnxembeddingtofile(service, text)

	//testqdrantinsert(context.Background(), service, ENV.QdrantMovieCollection)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", ENV.Embeddingserviceport))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcserver := grpc.NewServer()
	moviepb.RegisterMovieServiceServer(grpcserver, servercontroller)
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcserver.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
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
