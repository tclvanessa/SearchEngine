package main

import (
	"log"
	"strings"
	"sync"
	"unicode"

	"github.com/kljensen/snowball"
)

var mu sync.RWMutex

// Function to index the words into the db
func (mm *Indices) indexDB(inC chan ExtractResult) {
	// defer wg.Done()
	for ex := range inC {
		mm.totalDoc++
		mu.Lock()

		mm.updateTitle(ex.url, ex.title) // Insert title into db

		for _, sentence := range ex.sentences {
			wordsInSentence := make([]string, 0)
			stemmedNoStopWords := make([]string, 0)

			// fmt.Println(sentence)
			sentenceID, err := mm.insertSentence(sentence, ex.url) // Insert sentence into db
			if err != nil {
				log.Fatal("Error inserting sentences: ", err)
			}

			// Make slice of words out of sentence
			wordsInSentence = strings.FieldsFunc(sentence, func(r rune) bool {
				return !unicode.IsLetter(r) && !unicode.IsNumber(r)
			})

			// Stemming and stop words
			for _, word := range wordsInSentence {
				stem, err := snowball.Stem(word, "english", true)
				if err != nil {
					log.Println("Failed to stem", word, err)
				}
				// Check if the stemmed word is not a stop word
				if _, exists := mm.stopWordsMap[stem]; !exists {
					stemmedNoStopWords = append(stemmedNoStopWords, stem)
				}
			}

			for _, w := range stemmedNoStopWords {
				mm.insertTerm(w)                      // Insert stemmed word into db
				mm.insertCount(ex.url, w, sentenceID) // Insert hits count into db
			}

			// Insert for bigrams
			for i := 0; i < len(stemmedNoStopWords)-1; i++ {
				// fmt.Println(stemmedNoStopWords[i] + " " + stemmedNoStopWords[i+1])
				mm.insertBigram(stemmedNoStopWords[i], stemmedNoStopWords[i+1], ex.url, sentenceID)
			}
		}

		mu.Unlock()
	}
}
