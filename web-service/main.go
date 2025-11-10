package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	moviepb "github.com/Zpuspokusumo/kino-vectors/contract/golang/movie-services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PageData struct {
	Films           []*moviepb.MovieInfo
	AvailableGenres []string
}

type FilmService struct {
	client moviepb.MovieServiceClient
}

func setupservice() (*FilmService, error) {
	conn, err := grpc.NewClient("localhost:1555", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := moviepb.NewMovieServiceClient(conn)
	return &FilmService{
		client: client,
	}, nil
}

func (service *FilmService) SearchMoviesWithTopMovies(ctx context.Context, req *moviepb.RecommendMoviesRequest) (*moviepb.MovieInfos, error) {
	movielist, err := service.client.RecommendMovies(ctx, req)
	// add specific movies in genres here
	newmovies := &moviepb.MovieInfos{
		Movies: movielist.Movies,
	}
	return newmovies, err
}

func (service *FilmService) SearchMoviesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//get from query params
		queryParams := r.URL.Query()
		query := queryParams.Get("query")
		genreslist := r.URL.Query()["genres"]
		//pagination
		//qty := queryParams.Get("qty")
		//offset := queryParams.Get("offset")

		movierequest := &moviepb.RecommendMoviesRequest{
			TextQuery: query,
			Genres:    genreslist,
			Quantity:  0,
		}

		fmt.Printf("%+v\n\n", movierequest)

		_ = movierequest
		movielist, err := service.SearchMoviesWithTopMovies(r.Context(), movierequest)
		if err != nil {
			handleError(w, r, http.StatusBadRequest, "error fetching movies", "")
		}
		//movielist, _ := readdummycsv()

		//write response
		films := map[string][]*moviepb.MovieInfo{
			"Films": movielist.Movies,
		}
		tmpl := template.Must(template.ParseFiles("film-list-element.html"))
		tmpl.Execute(w, films)
	}
}

func readdummycsv() (*moviepb.MovieInfos, error) {
	file, err := os.Open("movies_reduced.csv") // Replace "data.csv" with your file name
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV: %v", err)
	}

	movies := []*moviepb.MovieInfo{}
	for i, record := range records {
		if i == 0 {
			continue
		}
		fmt.Printf("Row %d: %v\n", i, record)
		id, titleyear, genres := record[0], record[1], record[2]

		lentitle := len([]rune(titleyear))
		title := string([]rune(titleyear)[:lentitle-6])
		year, err := strconv.Atoi(string([]rune(titleyear)[lentitle-5 : lentitle-1]))
		if err != nil {
			fmt.Println(err)
			year = 0
		}

		genreslist := strings.Split(genres, "|")
		movie := &moviepb.MovieInfo{
			Id:    id,
			Title: string(title),
			Year:  uint32(year),
			Genre: genreslist,
		}
		movies = append(movies, movie)
	}
	return &moviepb.MovieInfos{
		Movies: movies,
	}, err
}

func GetDummyGenres() []string {
	return []string{"Sci-Fi", "Action", "Horror", "Romance"}
}

func main() {
	fmt.Println("Go app...")

	service, err := setupservice()
	if err != nil {
		log.Fatalf("fail to setup: %v", err)
	}

	genres := GetDummyGenres()

	// handler function #1 - returns the index.html template, with film data
	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		data := PageData{
			Films: []*moviepb.MovieInfo{
				{Title: "The Godfather", Director: "Francis Ford Coppola"},
				{Title: "Blade Runner", Director: "Ridley Scott"},
				{Title: "The Thing", Director: "John Carpenter"},
			},
			AvailableGenres: genres,
		}
		tmpl.Execute(w, data)
	}

	// define handlers
	http.HandleFunc("/", h1)
	http.HandleFunc("/fetch-film/", service.SearchMoviesHandler())

	// kamen rider faiz auto vajin code
	log.Fatal(http.ListenAndServe(":5821", nil))

}
