package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/kljensen/snowball"
)

type Result struct {
	Term, Title, URL, Sentence string
	Score                      float64
	Error                      bool
	ErrorMessage               string
}

type Results struct {
	Results []Result
}

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
	http.Handle("/search", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		term := r.URL.Query().Get("term")

		// Template
		t, err := template.ParseFiles("static/template.html")
		if err != nil {
			http.Error(w, "ParseFiles Error", http.StatusInternalServerError)
		}

		var searchResults []Result
		// If search query is bigram
		if mm.isTwoWords(term) {
			// Splitting term into two words
			terms := strings.Fields(term)
			firstTerm, secondTerm := terms[0], terms[1]
			if secondTerm == "authors" {
				secondTerm = "author"
			}
			// Searching bigram tfidf and getting hits
			hits, err := mm.searchBigramTfIdf(firstTerm, secondTerm)

			if len(hits) == 0 || err != nil {
				searchResults = append(searchResults, Result{
					Error:        true,
					ErrorMessage: "Word not found.",
				})
			} else {
				for _, h := range hits {
					title := mm.getTitle(h.URL)                                    // Get title from db using url
					sentence := mm.getSentenceBigram(firstTerm, secondTerm, h.URL) // Get sentence for bigrams
					searchResults = append(searchResults, Result{
						Term:     term,
						Title:    title,
						URL:      h.URL,
						Sentence: sentence,
						Score:    h.Score,
					})
				}
			}
		} else {
			// If wildcard box result empty, do regular search
			// If not, do wildcard search for the term
			if r.URL.Query().Get("wildcard") == "" {
				hits, err := mm.searchTfIdfDB(term)
				if len(hits) == 0 || err != nil {
					searchResults = append(searchResults, Result{
						Error:        true,
						ErrorMessage: "Word not found.",
					})
				} else {
					for _, h := range hits {
						title := mm.getTitle(h.URL)             // Get title from db using url
						sentence := mm.getSentence(term, h.URL) // Get sentence
						searchResults = append(searchResults, Result{
							Term:     term,
							Title:    title,
							URL:      h.URL,
							Sentence: sentence,
							Score:    h.Score,
						})
					}
				}
			} else {
				// Wildcard
				hitsWithTerms, err := mm.searchWildcardTfIdf(term)
				if len(hitsWithTerms) == 0 || err != nil {
					searchResults = append(searchResults, Result{
						Error:        true,
						ErrorMessage: "Word not found.",
					})
				} else {
					for _, h := range hitsWithTerms {
						title := mm.getTitle(h.URL) // Get title from db using url
						searchResults = append(searchResults, Result{
							Term:     h.Term,
							Title:    title,
							URL:      h.URL,
							Sentence: mm.getSentence(h.Term, h.URL),
							Score:    h.Score,
						})
					}
				}
			}
		}
		// Put searchResults into the results to send to execute
		results := Results{
			Results: searchResults,
		}

		// Execute results using template
		err = t.Execute(w, results)
		if err != nil {
			http.Error(w, "Execute Error", http.StatusInternalServerError)
		}
	}))
	go http.ListenAndServe(":8080", nil)
}

func (mm *Indices) searchTfIdfDB(term string) (Hits, error) {
	return mm.TfIdfDB(mm.searchDB(term))
}

func (mm *Indices) searchWildcardTfIdf(term string) (HitsWithTerms, error) {
	return mm.WildcardTfIdfDB(mm.searchDB(term))
}

func (mm *Indices) searchBigramTfIdf(first, second string) (Hits, error) {
	return mm.BigramTfIdfDB(mm.searchDB(first), mm.searchDB(second))
}

func (mm *Indices) isTwoWords(query string) bool {
	// Simple check if the query contains a space character
	return strings.Contains(query, " ")
}
