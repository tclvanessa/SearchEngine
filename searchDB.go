package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/kljensen/snowball"
)

func (mm *Indices) searchDB(term string) string {
	if stem, err := snowball.Stem(term, "english", true); err == nil {
		// fmt.Println(stem)
		return stem
	} else {
		log.Println("Failed to stem: " + term)
	}
	return ""
}

func serveDB(mm *Indices) {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/top10/", http.StripPrefix("/top10/", http.FileServer(http.Dir("top10"))))
	http.Handle("/search", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		term := r.URL.Query().Get("term")

		// If search query is bigram
		if mm.isTwoWords(term) {
			hits, err := mm.searchBigramTfIdf(term)
			if len(hits) == 0 || err != nil {
				io.WriteString(w, "Not found: "+term)
			} else {
				io.WriteString(w, "Number of matches for: "+term+"\n")
				for _, h := range hits {
					title := mm.getTitle(h.URL) // Get title from db using url
					line := fmt.Sprintf("%s\n%s: %v\n\n", title, h.URL, h.Score)
					io.WriteString(w, line)
				}
			}
		} else {
			// If wildcard box result empty, do regular search
			// If not, do wildcard search for the term
			if r.URL.Query().Get("wildcard") == "" {
				hits, err := mm.searchTfIdfDB(term)
				if len(hits) == 0 || err != nil {
					io.WriteString(w, "Not found: "+term)
				} else {
					io.WriteString(w, "Number of matches for: "+term+"\n")
					for _, h := range hits {
						title := mm.getTitle(h.URL) // Get title from db using url
						line := fmt.Sprintf("%s\n%s: %v\n\n", title, h.URL, h.Score)
						io.WriteString(w, line)
					}
				}
			} else {
				hits, err := mm.searchWildcardTfIdf(term)
				if len(hits) == 0 || err != nil {
					io.WriteString(w, "Not found: "+term)
				} else {
					io.WriteString(w, "Number of matches for: "+term+"\n")
					for _, h := range hits {
						title := mm.getTitle(h.URL) // Get title from db using url
						line := fmt.Sprintf("%s\n%s: %v\n\n", title, h.URL, h.Score)
						io.WriteString(w, line)
					}
				}
			}
		}
	}))
	go http.ListenAndServe(":8080", nil)
}

func (mm *Indices) searchTfIdfDB(term string) (Hits, error) {
	return mm.TfIdfDB(mm.searchDB(term))
}

func (mm *Indices) searchWildcardTfIdf(term string) (Hits, error) {
	return mm.WildcardTfIdfDB(mm.searchDB(term))
}

func (mm *Indices) searchBigramTfIdf(term string) (Hits, error) {
	return mm.BigramTfIdfDB(mm.searchDB(term))
}

func (mm *Indices) isTwoWords(query string) bool {
	// Simple check if the query contains a space character
	return strings.Contains(query, " ")
}
