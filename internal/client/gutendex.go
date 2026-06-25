package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

var (
	ErrNoMatch = errors.New("no id match")
)

type Author struct {
	Name      string `json:"name"`
	BirthYear int    `json:"birth_year"`
	DeathYear int    `json:"death_year"` //DeathYear = 0 means the author is not dead yet.
}

type Format struct {
	ImageJpeg string `json:"image/jpeg"`
}

type Book struct {
	Id      int      `json:"id"`
	Title   string   `json:"title"`
	Authors []Author `json:"authors"`
	Formats Format   `json:"formats"`
}

type Gutendex struct {
	Count   int    `json:"count"`
	Results []Book `json:"results"`
}

// SearchBooksByString uses the passed parameter as the query parameter for gutendex API use.
// It returns a slice of Books, built with the results from the search API.
// It accepts an empty string, since Gutendex API accepts it.

// Error returned as it is, since it's handled by the search books handler.
func SearchBooksByString(s string) ([]Book, error) {

	// Parse URL
	url := "https://gutendex.com/books/?search=" + url.QueryEscape(s)

	// Create http Client
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response body, save it into Results []Books and return it.
	// g.Results returns 32 Books at most, Gutendex rule.
	var g Gutendex

	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, err
	}
	return g.Results, nil
}

// SearchBooksByID queries Gutendex API to check if a Book exists given a specific, single ID.
// This is used for fast checking if a POST /review request body matches any Book in the Gutendex database.

// Returns an  if there is no match.
// Chose returning an error over a bool to differentiate No Match from other errors.
func SearchBooksByID(id int) error {

	strID := strconv.Itoa(id)
	// Parse URL
	url := "https://gutendex.com/books/?ids=" + url.QueryEscape(strID)

	// Create http Client
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response body, save it into Results []Books and return it.
	// g.Results returns 32 Books at most, Gutendex rule.
	var g Gutendex

	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return err
	}

	// len(slice) == 0 better than slice == nil.
	// It covers both nil and empty slice.
	//
	if len(g.Results) == 0 {
		return ErrNoMatch
	}
	return nil
}
