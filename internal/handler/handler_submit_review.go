package handler

import (
	"books/internal/gutendex"
	"books/internal/model"
	"books/internal/service"
	"books/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ReviewHandler struct {
	service *service.ReviewService
}

func NewReviewHandler(s *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		service: s,
	}
}

func (h *ReviewHandler) SubmitReview(w http.ResponseWriter, r *http.Request) {

	// Bad request if review is not a string, id and scores are not int?
	var req model.ReviewCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to decode response body", err)
		return
	}

	// If the book exists and the review body is validated, generate UUID for the reviewjob.
	reviewID, err := h.service.SubmitReview(r.Context(), req)
	if err != nil {
		if errors.Is(err, gutendex.ErrNoMatch) {
			utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("no match for the %d Book ID", req.BookID), err)
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "failed to submit review", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusAccepted, map[string]string{
		"review_id": reviewID,
	})

}
