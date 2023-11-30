package main

import (
	"log"
	"sort"
	"strings"
)

// Sort in descending order
func (h Hits) Less(i, j int) bool {
	// If score isn't the same, sort by score
	if h[i].Score != h[j].Score {
		return h[i].Score > h[j].Score
	} else {
		// Sort by URL if the score is the same
		return h[i].URL > h[j].URL
	}
}

// Swapping function for sorting hits
func (h Hits) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Returns the Hits for a term
func (mm *Indices) TfIdfDB(term string) (Hits, error) {
	hits := Hits{}

	// Using helper function to get the term id from db
	termID, _, err := mm.getTermID(term)
	if err != nil {
		log.Printf("Error in TfIdf getting termID: %v\n", err)
		return hits, err
	}

	// Using helper function to get a slice of url ids from a term id
	urlIDsForWord := mm.getURLIDsFromTermID(termID)

	// Range over that slice of url ids and get the url itself
	// to calculate the tfidf score and append the url and score to hits
	for _, urlID := range urlIDsForWord {
		url := mm.getURL(urlID)
		score := mm.getTFIDF(term, url)
		hits = append(hits, Hit{
			url,
			score,
		})
	}
	sort.Sort(hits) // Sort hits
	return hits, nil
}

// Returns Hits for wildcard search
func (mm *Indices) WildcardTfIdfDB(term string) (Hits, error) {
	hits := Hits{}

	// Getting rows of wildcard term ids
	termIDs, err := mm.getTermIDsLike(term)
	if err != nil {
		log.Printf("Error in WildcardTfIdfDB getting termID: %v\n", err)
		return hits, err
	}

	// Use helper function to get a slice of url ids from the rows of term ids
	for _, termID := range termIDs {
		// Using helper function to get a slice of URL IDs from a term ID
		urlIDsForWord := mm.getURLIDsFromTermID(termID)

		// Range over that slice of URL IDs and get the URL itself
		// to calculate the TF-IDF score and append the URL and score to hits
		for _, urlID := range urlIDsForWord {
			termFromID := mm.getTerm(termID)
			url := mm.getURL(urlID)
			score := mm.getTFIDF(termFromID, url)
			hits = append(hits, Hit{
				url,
				score,
			})
		}
	}
	sort.Sort(hits) // Sort hits
	return hits, nil
}

// Returns the Hits for bigram term
func (mm *Indices) BigramTfIdfDB(term string) (Hits, error) {
	hits := Hits{}

	// Splitting bigram into two words
	terms := strings.Fields(term)
	firstTerm, secondTerm := terms[0], terms[1]
	// fmt.Println(firstTerm + " & " + secondTerm)

	// Using helper function to get the term id from db
	firstTermID, _, err := mm.getTermID(firstTerm)
	if err != nil {
		log.Printf("Error in BigramTfIdfDB getting termID1: %v\n", err)
		return hits, err
	}
	secondTermID, _, err := mm.getTermID(secondTerm)
	if err != nil {
		log.Printf("Error in BigramTfIdfDB getting termID2: %v\n", err)
		return hits, err
	}

	// Using helper function to get a slice of url ids from two term ids
	urlIDsForWord := mm.getURLIDsFromTwoTermIDs(firstTermID, secondTermID)

	// Range over that slice of url ids and get the url itself
	// to calculate the tfidf score and append the url and score to hits
	for _, urlID := range urlIDsForWord {
		url := mm.getURL(urlID)
		score := mm.getBigramTFIDF(firstTerm, secondTerm, url)
		hits = append(hits, Hit{
			url,
			score,
		})
	}
	sort.Sort(hits) // Sort hits
	return hits, nil
}

// Calculating the TFIDF score here
func (mm *Indices) getTFIDF(term string, url string) float64 {
	// Use helper func to get the count of the term in the doc
	tc, err := mm.getCount(url, term)
	if err != nil {
		log.Println("Error getting term count: ", err)
	}

	// Use helper func to get the total words in the doc
	wordsInDoc, err := mm.getWordsInDoc(url)
	if err != nil {
		log.Println("Error getting words in doc: ", err)
	}

	// Calculate the term frequency by dividing the term count by the total number of words in the doc
	tf := float64(tc) / float64(wordsInDoc)

	// Use helper func to get the total num of docs with the term in it
	docTerm, err := mm.getTotalDocWTerm(term)
	if err != nil {
		log.Println("Error getting num of documents with term: ", err)
	}

	// Calculate df by dividing total docs with term by total num of docs in general
	df := float64(docTerm) / float64(mm.totalDoc)
	idf := 1.0 / df   // Calculate idf
	score := tf * idf // Get the score by multiplying tf by idf

	return score
}

// Calculate Tf-Idf score with bigram term
func (mm *Indices) getBigramTFIDF(first, second, url string) float64 {
	// Use helper func to get the count of the bigram term in the doc
	tc, err := mm.getBigramCount(url, first, second)
	if err != nil {
		log.Println("Error getting bigram term count: ", err)
	}

	// Use helper func to get the total words in the doc
	wordsInDoc, err := mm.getWordsInDocBigram(url)
	if err != nil {
		log.Println("Error getting words in doc bigram: ", err)
	}

	// Calculate the term frequency by dividing the term count by the total number of words in the doc
	tf := float64(tc) / float64(wordsInDoc)

	// Use helper func to get the total num of docs with the term in it
	docTerm, err := mm.getTotalDocWTermBigram(first, second)
	if err != nil {
		log.Println("Error getting num of documents with bigram term: ", err)
	}

	// Calculate df by dividing total docs with term by total num of docs in general
	df := float64(docTerm) / float64(mm.totalDoc)
	idf := 1.0 / df   // Calculate idf
	score := tf * idf // Get the score by multiplying tf by idf

	return score
}
