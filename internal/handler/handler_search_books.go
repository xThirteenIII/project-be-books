package handler

import (
	"books/internal/gutendex"
	"books/internal/utils"
	"net/http"
)

// SearchBooksByKeywords queries Gutendex GET API to search books given a string.
// Don't need a service for this since the logic is small.
func SearchBooksByKeywords(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	// Not checking for empty q since Gutendex API allows it.
	// It returns first 32 books.

	books, err := gutendex.SearchBooksByString(q)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch books", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, books)
}
