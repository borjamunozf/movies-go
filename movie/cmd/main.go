package main

import (
	"log"
	movie "microgomovies/movie/internal/controller"
	metadatagateway "microgomovies/movie/internal/gateway/metadata/http"
	ratinggateway "microgomovies/movie/internal/gateway/rating/http"
	httphandler "microgomovies/movie/internal/handler/http"
	"net/http"
)

func main() {
	metadataGateway := metadatagateway.New("localhost:8081")
	ratingGateway := ratinggateway.New("localhost:8082")

	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))

	log.Println("Starting Movie service")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
