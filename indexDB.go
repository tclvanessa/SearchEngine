package main

import "sync"

var mu sync.RWMutex

// Function to index the words into the db
func (mm *Indices) indexDB(inC chan ExtractResult) {
	// defer wg.Done()
	for ex := range inC {
		mm.totalDoc++
		mu.Lock()

		mm.updateTitle(ex.url, ex.title) // Insert title into db

		// Range over words
		for _, w := range ex.words {
			mm.insertTerm(w)          // Insert word into db
			mm.insertCount(ex.url, w) // Insert hits count into db

		}
		for i := 0; i < len(ex.words)-1; i++ {
			// fmt.Println(ex.words[i] + " " + ex.words[i+1])
			mm.insertBigram(ex.words[i], ex.words[i+1], ex.url)
		}

		mu.Unlock()
	}
}
