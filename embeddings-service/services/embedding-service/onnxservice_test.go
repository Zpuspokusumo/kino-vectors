package embeddingservice

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestEmbedding(t *testing.T) {
	data, err := Onnxservice("text here")
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
