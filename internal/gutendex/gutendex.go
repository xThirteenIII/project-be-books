// package client queries Gutendex APIs
package gutendex

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
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

type GutendexClient struct{}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// SearchBooksByString uses the passed parameter as the query parameter for gutendex API use.
// It returns a slice of Books, built with the results from the search API.
// It accepts an empty string, since Gutendex API accepts it.

// Error returned as it is, since it's handled by the search books handler.
func SearchBooksByString(s string) ([]Book, error) {

	endpoint := "https://gutendex.com/books/?search=" + url.QueryEscape(s)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("gutendex search returned status %d", resp.StatusCode)
	}

	// Decode response body, save it into Results []Books and return it.
	// g.Results returns 32 Books at most, Gutendex rule.
	var g Gutendex

	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, err
	}
	return g.Results, nil
}

// SearchByID wraps SearchBooksByID, mainly for testing purposes.
// GutendexClient implements BookLookup interface.
func (c *GutendexClient) SearchByID(id int) (Book, error) {
	return SearchBooksByID(id)
}

// SearchBooksByID queries Gutendex API to check if a Book exists given a specific, single ID.
// This is used for fast checking if a POST /review request body matches any Book in the Gutendex database.
func SearchBooksByID(id int) (Book, error) {

	strID := strconv.Itoa(id)
	endpoint := "https://gutendex.com/books/?ids=" + url.QueryEscape(strID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return Book{}, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return Book{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return Book{}, fmt.Errorf("gutendex lookup returned status %d", resp.StatusCode)
	}

	// Decode response body, save it into Results []Books and return it.
	// g.Results returns 32 Books at most, Gutendex rule.
	var g Gutendex

	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return Book{}, err
	}

	// len(slice) == 0 better than slice == nil.
	// It covers both nil and empty slice.
	if len(g.Results) == 0 {
		return Book{}, ErrNoMatch
	}
	return g.Results[0], nil
}
