package handler

import (
	"books/internal/gutendex"
	"books/internal/model"
	"books/internal/repo"
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

	var req model.ReviewCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to decode response body", err)
		return
	}

	// If the book exists and the review body is validated, generate UUID for the reviewjob.
	reviewID, err := h.service.SubmitReview(r.Context(), req)
	if err != nil {
		if errors.Is(err, gutendex.ErrNoMatch) {
			utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("no match for the %d Book ID", req.BookID), err)
			return
		}

		if errors.Is(err, service.ErrValidation) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error(), err)
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "failed to submit review", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusAccepted, map[string]string{
		"review_id": reviewID,
	})

}

func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	reviewID := r.PathValue("id")
	if reviewID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing review id", nil)
		return
	}

	review, err := h.service.GetReview(r.Context(), reviewID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			utils.RespondWithError(w, http.StatusNotFound, "review not found", err)
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "failed to get review", err)
		return
	}

	// Answer 202 if status is pending.
	if review.Status == model.StatusPending {
		utils.RespondWithJSON(w, http.StatusAccepted, review)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, review)
}

func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	reviewID := r.PathValue("id")
	if reviewID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing review id", nil)
		return
	}

	var req model.ReviewUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to decode request body", err)
		return
	}

	if err := h.service.UpdateReview(r.Context(), reviewID, req); err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			utils.RespondWithError(w, http.StatusNotFound, "review not found", err)
			return
		}

		if errors.Is(err, service.ErrValidation) {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error(), err)
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "failed to update review", err)
		return
	}

	// GET /review/{id}
	// No need to return a response with a review body. It would mean calling GetByID again.
	w.WriteHeader(http.StatusNoContent)
}

func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	reviewID := r.PathValue("id")
	if reviewID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing review id", nil)
		return
	}

	if err := h.service.DeleteReview(r.Context(), reviewID); err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			utils.RespondWithError(w, http.StatusNotFound, "review not found", err)
			return
		}

		utils.RespondWithError(w, http.StatusInternalServerError, "failed to delete review", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
