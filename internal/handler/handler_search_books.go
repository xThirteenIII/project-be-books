package handler

import (
	"books/internal/client"
	"books/internal/utils"
	"net/http"
)

func SearchBooksByKeywords(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	// Not checking for empty q since Gutendex method allows it.
	// It returns first 32 books.

	books, err := client.SearchBooksByString(q)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch books", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, books)
}
