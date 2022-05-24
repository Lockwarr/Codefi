package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Lockwarr/codefi/pkg/scraper"
	"github.com/Lockwarr/codefi/services/links/domain"
	"github.com/Lockwarr/codefi/services/links/handler"
	"github.com/Lockwarr/codefi/services/links/repository"
	"github.com/go-chi/chi/v5"
)

var port = ":8080" // could be moved to cfg

func main() {
	log.Println("Starting links service")
	router := chi.NewRouter()
	repo := repository.NewInMemoryDB()
	scraper := scraper.NewScraper()
	linksProcessor := domain.NewLinksProcessor(repo, scraper)
	h := handler.NewHandler(linksProcessor)

	router.Route("/api/v1/", func(r chi.Router) {
		r.Post("/links", h.ProcessBatch)
		r.Get("/links/{batchID}", h.GetBatch)
	})

	if err := http.ListenAndServe(port, router); err != nil {
		log.Println(err.Error(), "failed to start http server")
		os.Exit(1)
	}
}
