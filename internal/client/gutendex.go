package client

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Author struct {
	Name      string `json:"name"`
	BirthYear string `json:"birth_year"`
	DeathYear string `json:"death_year"`
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

func SearchBooksByString(s string) ([]Book, error) {
	url := "https://gutendex.com/books/?search=" + url.QueryEscape(s)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// g.Results returns 32 Books at most, Gutendex rule.
	var g Gutendex

	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, err
	}
	return g.Results, nil
}
