package httphandler

import (
	"encoding/json"
	"errors"
	"log"
	rating "microgomovies/rating/internal/controller"
	model "microgomovies/rating/pkg"
	"net/http"
	"strconv"
)

type RatingHandler struct {
	c *rating.Controller
}

func New(c *rating.Controller) *RatingHandler {
	return &RatingHandler{c}
}

func (h *RatingHandler) Handle(w http.ResponseWriter, req *http.Request) {
	recordID := model.RecordID(req.FormValue("id"))

	if recordID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recordType := model.RecordType(req.FormValue("type"))
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodGet:
		r, err := h.c.GetAggregatedRating(req.Context(), recordID, recordType)
		if err != nil && errors.Is(err, rating.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(r); err != nil {
			log.Printf("response encode errors: %v\n", err)
		}
	case http.MethodPut:
		userID := model.UserID(req.FormValue("userId"))
		r, err := strconv.ParseFloat(req.FormValue("value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		if err := h.c.PutRating(req.Context(), recordID, recordType, &model.Rating{UserID: userID, Value: model.RatingValue(r)}); err != nil {
			log.Printf("Saving rating error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
