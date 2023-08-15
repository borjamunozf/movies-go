package main

import (
	"log"
	metadata "microgomovies/metadata/internal/controller"
	httphandler "microgomovies/metadata/internal/handler/http"
	"microgomovies/metadata/internal/repository/memory"
	"net/http"
)

func main() {
	log.Println("Starting movie metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)

	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
