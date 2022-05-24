// Package handler contains Handler logic for the links service
package handler

import (
	"errors"
	"net/http"

	"github.com/Lockwarr/codefi/pkg/helpers"
	"github.com/Lockwarr/codefi/services/links"
	"github.com/Lockwarr/codefi/services/links/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var ErrNoUrlsForProcessing = errors.New("no urls for processing")
var ErrRetrievingFile = errors.New("bad file")

type Handler struct {
	linksProcessor links.Processor
}

// NewHandler ..
func NewHandler(linksProcessor links.Processor) *Handler {
	return &Handler{linksProcessor: linksProcessor}
}

// StartBatchProcessing - handler to start processing of batch of urls
// passed in a file with multi-line text with valid url on each line.
func (h *Handler) ProcessBatch(w http.ResponseWriter, r *http.Request) {
	// FormFile returns the first file for the given key `urlsFile`
	file, _, err := r.FormFile("urlsFile")

	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, links.Response{Errors: []string{ErrRetrievingFile.Error()}})
		return
	}
	defer file.Close()

	urls, err := helpers.GatherUrls(file)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, links.Response{Errors: []string{err.Error()}})
		return
	}

	if len(urls) == 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, links.Response{Errors: []string{ErrNoUrlsForProcessing.Error()}})
		return
	}

	results, err := h.linksProcessor.ProcessBatch(r.Context(), links.ProcessBatchRequest{URLs: urls})
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, links.Response{Errors: []string{links.ErrInternalServerError.Error()}})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, links.Response{Data: links.ProcessBatchResponse{Results: results}})
}

// GetBatch - handler for getting processed links by batch ID
func (h *Handler) GetBatch(w http.ResponseWriter, r *http.Request) {
	batchID := chi.URLParam(r, "batchID")

	results, err := h.linksProcessor.GetBatch(r.Context(), links.GetBatchRequest{BatchID: batchID})
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrBatchNotFound): // batch not found
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, links.Response{Errors: []string{repository.ErrBatchNotFound.Error()}})
			return
		default: // generic response to not leak details for all other errors
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, links.Response{Errors: []string{links.ErrInternalServerError.Error()}})
			return
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, links.Response{Data: links.GetBatchResponse{Results: results}})
}
