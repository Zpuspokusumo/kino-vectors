package controller

import (
	"context"
	embeddingservice "kino-vectors/services/embedding-service"

	moviepb "github.com/Zpuspokusumo/kino-vectors/contract/golang/movie-services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Controller struct {
	moviepb.UnimplementedMovieServiceServer
	service *embeddingservice.EmbeddingServiceONNX
}

// NewMyServer creates and returns a new server instance.
func New(service *embeddingservice.EmbeddingServiceONNX) *Controller {
	return &Controller{
		service: service,
	}
}

// --- Service Method Implementations ---
func (controller *Controller) ProcessMovie(context.Context, *moviepb.MovieInfo) (*moviepb.ProcessMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessMovie not implemented")
}
func (Controller *Controller) ProcessMovies(context.Context, *moviepb.MovieInfos) (*moviepb.ProcessMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessMovies not implemented")
}
func (Controller *Controller) RecommendMovies(context.Context, *moviepb.RecommendMoviesRequest) (*moviepb.RecommendMoviesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecommendMovies not implemented")
}
