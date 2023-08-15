package main

import (
	rating "microgomovies/rating/internal/controller"
	httphandler "microgomovies/rating/internal/handler/http"
	"microgomovies/rating/internal/repository/memory"
	"net/http"
)

func main() {
	repo := memory.New()
	ctrl := rating.New(repo)
	handler := httphandler.New(ctrl)

	http.Handle("/rating", http.HandlerFunc(handler.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
