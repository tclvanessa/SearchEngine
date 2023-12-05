package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Function to create the database with tables
func (mm *Indices) createDB() {
	os.Remove("sqlindex.db")
	// Open a connection to the SQLite database (or create a new one if it doesn't exist)
	db, err := sql.Open("sqlite3", "sqlindex.db")
	if err != nil {
		log.Fatal("Error in opening connection to db: ", err)
	}

	// Create the 'URLs' table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS urls (
        id INTEGER NOT NULL PRIMARY KEY,
        name TEXT,
		title TEXT
    )`)
	if err != nil {
		log.Fatal("Error in creating url table: ", err)
	}

	// Create the 'terms' table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS terms (
        id INTEGER NOT NULL PRIMARY KEY,
        name TEXT
    )`)
	if err != nil {
		log.Fatal("Error in creating terms table: ", err)
	}

	// Create the 'hits' table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS hits (
        id INTEGER NOT NULL PRIMARY KEY,
		urlID INTEGER,
		termID INTEGER,
		sentenceID INTEGER,
		count INTEGER,
        FOREIGN KEY(termID) REFERENCES terms(id),
		FOREIGN KEY(urlID) REFERENCES urls(id),
		FOREIGN KEY(sentenceID) REFERENCES sentences(id)
    )`)
	if err != nil {
		log.Fatal("Error in creating hits table: ", err)
	}

	// Create the 'bigram' table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS bigram (
        id INTEGER NOT NULL PRIMARY KEY,
        termID1 TEXT,
		termID2 TEXT,
		urlID INTEGER,
		sentenceID INTEGER,
		count INTEGER,
		FOREIGN KEY(termID1) REFERENCES terms(id),
		FOREIGN KEY(termID2) REFERENCES terms(id),
		FOREIGN KEY(urlID) REFERENCES urls(id),
		FOREIGN KEY(sentenceID) REFERENCES sentences(id)
    )`)
	if err != nil {
		log.Fatal("Error in creating bigram table: ", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sentences (
        id INTEGER NOT NULL PRIMARY KEY,
        name TEXT,
		urlID INTEGER,
		FOREIGN KEY(urlID) REFERENCES urls(id)
    )`)
	if err != nil {
		log.Fatal("Error in creating sentences table: ", err)
	}

	// Prepared statements
	mm.insTermStmt, err = db.Prepare("INSERT INTO terms(name) VALUES(?);")
	if err != nil {
		log.Fatal("insTermStmt: ", err)
	}

	mm.selTerm, err = db.Prepare("SELECT name FROM terms WHERE id=?;")
	if err != nil {
		log.Fatal("selTerm: ", err)
	}

	mm.selTermIDStmt, err = db.Prepare("SELECT id FROM terms WHERE name=?;")
	if err != nil {
		log.Fatal("selTermIDStmt: ", err)
	}

	mm.selTermIDLikeStmt, err = db.Prepare("SELECT id FROM terms WHERE name LIKE ?;")
	if err != nil {
		log.Fatal("selTermIDLikeStmt: ", err)
	}

	mm.selExisTermStmt, err = db.Prepare("SELECT EXISTS(SELECT 1 FROM terms WHERE name=?);")
	if err != nil {
		log.Fatal("selExisTermStmt: ", err)
	}

	mm.insURLStmt, err = db.Prepare("INSERT INTO urls(name) VALUES(?);")
	if err != nil {
		log.Fatal("insURLStmt: ", err)
	}

	mm.selURLIDStmt, err = db.Prepare("SELECT id FROM urls WHERE name=?;")
	if err != nil {
		log.Fatal("selURLIDStmt: ", err)
	}

	mm.selURLStmt, err = db.Prepare("SELECT name FROM urls WHERE id=?;")
	if err != nil {
		log.Fatal("selURLStmt: ", err)
	}

	mm.selExisURLStmt, err = db.Prepare("SELECT EXISTS(SELECT 1 FROM urls WHERE name=?);")
	if err != nil {
		log.Fatal("selExisURLStmt: ", err)
	}

	mm.selURLIDFromTermIDStmt, err = db.Prepare("SELECT urlID FROM hits WHERE termID=?;")
	if err != nil {
		log.Fatal("selURLIDFromTermIDStmt: ", err)
	}

	mm.selURLIDFromBigramStmt, err = db.Prepare("SELECT urlID FROM bigram WHERE termID1=? AND termID2=?;")
	if err != nil {
		log.Fatal("selURLIDFromBigramStmt: ", err)
	}

	mm.selCount, err = db.Prepare("SELECT count FROM hits WHERE urlID=? AND termID=? AND sentenceID=?;")
	if err != nil {
		log.Fatal("selCount: ", err)
	}

	mm.selCountTermID, err = db.Prepare("SELECT COUNT(*) FROM hits WHERE termID=?;")
	if err != nil {
		log.Fatal("selCountTermID: ", err)
	}

	mm.selCountTermIDBigram, err = db.Prepare("SELECT COUNT(*) FROM bigram WHERE termID1=? AND termID2=?;")
	if err != nil {
		log.Fatal("selCountTermIDBigram: ", err)
	}

	mm.updHits, err = db.Prepare("UPDATE hits SET count=? WHERE urlID=? AND termID=? AND sentenceID=?;")
	if err != nil {
		log.Fatal("updHits: ", err)
	}

	mm.insHits, err = db.Prepare("INSERT INTO hits (count, urlID, termID, sentenceID) VALUES (1, ?, ?, ?);")
	if err != nil {
		log.Fatal("insHits: ", err)
	}

	mm.selSumCount, err = db.Prepare("SELECT SUM(count) FROM hits WHERE urlID=?;")
	if err != nil {
		log.Fatal("selSumCount: ", err)
	}

	mm.selSumCountBigram, err = db.Prepare("SELECT SUM(count) FROM bigram WHERE urlID=?;")
	if err != nil {
		log.Fatal("selSumCount: ", err)
	}

	mm.insBigram, err = db.Prepare("INSERT INTO bigram (count, urlID, termID1, termID2, sentenceID) VALUES (1, ?, ?, ?, ?);")
	if err != nil {
		log.Fatal("insBigram: ", err)
	}

	mm.updBigram, err = db.Prepare("UPDATE bigram SET count=? WHERE urlID=? AND termID1=? AND termID2=? AND sentenceID=?;")
	if err != nil {
		log.Fatal("updBigram: ", err)
	}

	mm.selCountBigram, err = db.Prepare("SELECT count FROM bigram WHERE urlID=? AND termID1=? AND termID2=? AND sentenceID=?;")
	if err != nil {
		log.Fatal("selCountBigram: ", err)
	}

	mm.updTitle, err = db.Prepare("UPDATE urls SET title=? WHERE name=?;")
	if err != nil {
		log.Fatal("updTitle: ", err)
	}

	mm.selTitle, err = db.Prepare("SELECT title FROM urls WHERE name=?;")
	if err != nil {
		log.Fatal("selTitle: ", err)
	}

	mm.insSentenceStmt, err = db.Prepare("INSERT INTO sentences(name, urlID) VALUES(?, ?);")
	if err != nil {
		log.Fatal("insSentenceStmt: ", err)
	}

	mm.selSentence, err = db.Prepare("SELECT name FROM sentences WHERE id=?;")
	if err != nil {
		log.Fatal("selSentenceID: ", err)
	}

	mm.selSentenceBigram, err = db.Prepare("SELECT sentenceID FROM bigram WHERE urlID=? AND termID1=? AND termID2=?;")
	if err != nil {
		log.Fatal("selSentenceBigram: ", err)
	}

	mm.selSentenceID, err = db.Prepare("SELECT sentenceID FROM hits WHERE urlID=? AND termID=?;")
	if err != nil {
		log.Fatal("selSentenceID: ", err)
	}

	mm.selBigramSentence, err = db.Prepare("SELECT name FROM sentences WHERE id=?;")
	if err != nil {
		log.Fatal("selBigramSentence: ", err)
	}

	mm.selBigramSentenceID, err = db.Prepare("SELECT sentenceID FROM bigram WHERE urlID=? AND termID1=? AND termID2=?;")
	if err != nil {
		log.Fatal("selBigramSentenceID: ", err)
	}

	log.Println("Database and tables created successfully.")
	mm.db = db
}

// Insert sentence into sentences table
func (mm *Indices) insertSentence(sentence, url string) (int, error) {
	urlID, _, err := mm.getURLID(url)
	if err != nil {
		log.Println("Error getting url id in insertSentence for ", url)
	}

	// insert into table
	result, err1 := mm.insSentenceStmt.Exec(sentence, urlID)
	if err1 != nil {
		return 0, err1
	}
	lastInsertID, err2 := result.LastInsertId()
	if err2 != nil {
		return 0, err2
	}

	return int(lastInsertID), nil
}

// Get sentence for interface
func (mm *Indices) getSentence(term, url string) string {
	sentenceID, err := mm.getSentenceID(term, url)
	if err != nil {
		log.Println("Error getting sentence id in getSentence: ", err)
	}

	var sentence, nextSentence string
	for len(sentence) < 100 {
		err := mm.selSentence.QueryRow(sentenceID).Scan(&sentence)
		if err != nil {
			log.Println("Error getting sentence in getSentence: ", err)
			break // Break the loop if there's an error fetching the current sentence
		}

		err = mm.selSentence.QueryRow(sentenceID + 1).Scan(&nextSentence)
		if err != nil {
			log.Println("Error getting next sentence in getSentence: ", err)
			break // Break the loop if there's an error fetching the next sentence
		}

		sentence += " " + nextSentence
		sentenceID++
	}
	return sentence
}

// Get sentence for bigrams for interface
func (mm *Indices) getSentenceBigram(first, second, url string) string {
	bigramSentenceID, err := mm.getSentenceIDBigram(url, first, second)
	if err != nil {
		log.Println("Error getting bigram sentence id in getSentenceBigram: ", err)
	}

	var sentence, nextSentence string
	for len(sentence) < 100 {
		err := mm.selSentence.QueryRow(bigramSentenceID).Scan(&sentence)
		if err != nil {
			log.Println("Error getting sentence in getSentence: ", err)
			break // Break the loop if there's an error fetching the current sentence
		}

		err = mm.selSentence.QueryRow(bigramSentenceID + 1).Scan(&nextSentence)
		if err != nil {
			log.Println("Error getting next sentence in getSentence: ", err)
			break // Break the loop if there's an error fetching the next sentence
		}

		sentence += " " + nextSentence
		bigramSentenceID++
	}
	return sentence
}

func (mm *Indices) getSentenceID(term, url string) (int, error) {
	urlID, _, err1 := mm.getURLID(url)
	if err1 != nil {
		log.Println("Error getting url id in getSentenceID for: ", err1)
	}

	termID, _, err2 := mm.getTermID(term)
	if err2 != nil {
		log.Println("Error getting term id in getSentenceID: ", err2)
	}

	var sentenceID int
	err3 := mm.selSentenceID.QueryRow(urlID, termID).Scan(&sentenceID)
	if err3 != nil {
		log.Println("Error getting sentence id in getSentenceID: ", err3)
		return 0, err3
	}
	return sentenceID, nil
}

func (mm *Indices) getSentenceIDBigram(url, term1, term2 string) (int, error) {
	urlID, _, err1 := mm.getURLID(url)
	if err1 != nil {
		log.Println("Error getting url id in getSentenceIDBigram: ", err1)
		return -1, err1
	}
	firstTermID, _, err2 := mm.getTermID(term1)
	if err2 != nil {
		log.Println("Error getting term id 1 in getSentenceIDBigram: ", err2)
		return -1, err2
	}
	secondTermID, _, err3 := mm.getTermID(term2)
	if err3 != nil {
		log.Println("Error getting term id 2 in getSentenceIDBigram: ", err3)
		return -1, err3
	}

	var sentenceID int
	err := mm.selBigramSentenceID.QueryRow(urlID, firstTermID, secondTermID).Scan(&sentenceID)
	if err != nil {
		log.Println("Error selecting bigram sentence id: ", err)
		return 0, err
	}
	return sentenceID, nil
}

// Function for inserting title into urls table using update
func (mm *Indices) updateTitle(url string, title string) error {
	_, err := mm.updTitle.Exec(title, url)
	if err != nil {
		log.Fatal("Error updating title in urls table: ", err)
		return err
	}
	return nil
}

// Function to get title using a url
func (mm *Indices) getTitle(url string) string {
	var title string
	if err := mm.selTitle.QueryRow(url).Scan(&title); err == nil {
		return title
	} else {
		log.Fatal("Error getting title: ", err)
	}
	return ""
}

// Function for inserting URL into db
func (mm *Indices) insertURL(url string) (bool, error) {
	var exists bool
	// If url exists in db, return true
	if err := mm.selExisURLStmt.QueryRow(url).Scan(&exists); err != nil {
		return exists, err
	}

	// If url doesn't exist, insert into urls table
	if !exists {
		if _, err := mm.insURLStmt.Exec(url); err != nil {
			return exists, err
		}
	}
	return exists, nil
}

// Function for inserting term into db
func (mm *Indices) insertTerm(term string) error {
	var exists bool
	// If term exists in db, return error
	if err := mm.selExisTermStmt.QueryRow(term).Scan(&exists); err != nil {
		log.Println("QueryRow Error in insertTerm", err)
		return err
	}

	// If term doesn't exist, insert into terms table
	if !exists {
		if _, err := mm.insTermStmt.Exec(term); err != nil {
			log.Println("Error inserting in terms", err)
			return err
		}
	}
	return nil
}

// Function to get url id from the url
func (mm *Indices) getURLID(url string) (int, bool, error) {
	var id int
	if err := mm.selURLIDStmt.QueryRow(url).Scan(&id); err == nil {
		return id, true, nil
	} else {
		log.Fatalf("URL ID not found %v", err)
		return 0, false, err
	}
}

// Function to get the url from the url id
func (mm *Indices) getURL(urlID int) string {
	var url string
	if err := mm.selURLStmt.QueryRow(urlID).Scan(&url); err == nil {
		return url
	} else {
		log.Fatalf("URL not found: %v, %v", url, err)
		return ""
	}
}

// Function to get the term id from the term
func (mm *Indices) getTermID(term string) (int, bool, error) {
	var id int
	if err := mm.selTermIDStmt.QueryRow(term).Scan(&id); err == nil {
		return id, true, nil
	} else {
		log.Printf("Term ID not found for %v: %v", term, err)
		return 0, false, err
	}
}

// Function to get the term from the term id (Don't think I need this function, but just in case!)
func (mm *Indices) getTerm(termID int) string {
	var term string
	err := mm.selTerm.QueryRow(termID).Scan(&term)
	if err != nil {
		log.Fatalf("Term not found %v", err)
		return ""
	} else {
		return term
	}
}

// Function that returns a slice of url ids related to a term id
func (mm *Indices) getURLIDsFromTermID(termID int) []int {
	var urlsID []int

	// Querying for rows of url ids from a term id
	rows, err := mm.selURLIDFromTermIDStmt.Query(termID)
	if err != nil {
		log.Println("Error querying when getting urls ids of a term: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var urlID int
		err := rows.Scan(&urlID)
		if err != nil {
			log.Println("Error scanning rows for url ids in getURLIDsFromTermID: ", err)
		}
		urlsID = append(urlsID, urlID)
	}
	return urlsID
}

func (mm *Indices) getURLIDsFromTwoTermIDs(first, second int) []int {
	var urlsID []int

	// Querying for rows of url ids from two term ids
	rows, err := mm.selURLIDFromBigramStmt.Query(first, second)
	if err != nil {
		log.Println("Error querying when getting urls ids of two terms: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var urlID int
		err := rows.Scan(&urlID)
		if err != nil {
			log.Println("Error scanning rows for url ids in getURLIDsFromTwoTermIDs: ", err)
		}
		urlsID = append(urlsID, urlID)
	}
	return urlsID
}

// Helper function to get term IDs with wildcard searches
func (mm *Indices) getTermIDsLike(term string) ([]int, error) {
	var termIDs []int

	// Execute the LIKE query to get term IDs
	rows, err := mm.selTermIDLikeStmt.Query(term + "%")
	if err != nil {
		log.Println("Error querying when getting like term ids of a term: ", err)
		return termIDs, err
	}
	defer rows.Close()

	// Iterate over the rows and append term IDs to the slice
	for rows.Next() {
		var termID int
		err := rows.Scan(&termID)
		if err != nil {
			log.Fatal("Error scanning for term id in getTermIDsLike: ", err)
		}
		termIDs = append(termIDs, termID)
	}
	// fmt.Println(termIDs)

	return termIDs, nil
}

func (mm *Indices) insertBigram(first, second, url string, sentenceID int) error {
	// Get url and term ids
	urlID, _, err := mm.getURLID(url)
	if err != nil {
		log.Println("Error getting url id for ", url)
	}
	firstTermID, _, err := mm.getTermID(first)
	if err != nil {
		log.Println("Error getting term id for ", url)
	}
	secondTermID, _, err := mm.getTermID(second)
	if err != nil {
		log.Println("Error getting term id for ", url)
	}

	var hits int
	err = mm.selCountBigram.QueryRow(urlID, firstTermID, secondTermID, sentenceID).Scan(&hits)
	if err == nil {
		// If bigram exists, update bigram table
		hits++ // Incrementing number of hits

		_, err = mm.updBigram.Exec(hits, urlID, firstTermID, secondTermID, sentenceID)
		if err != nil {
			log.Fatalf("Error updating bigram table %v", err)
			return err
		}
		// fmt.Printf("Successfully updated bigram to %d for URL: %v and Term: %v\n", hits, url, term)
	} else {
		// If word doesn't exist, insert into bigram table
		if err == sql.ErrNoRows {
			hits = 1

			_, err = mm.insBigram.Exec(urlID, firstTermID, secondTermID, sentenceID)
			if err != nil {
				log.Fatalf("Error inserting into bigram table %v", err)
				return err
			}
		}
		// fmt.Printf("Successfully added to bigram - Hit: %d, URL: %v, Term: %v\n", hits, url, term)
	}
	return nil
}

// Function to get the term count of a term in a url
func (mm *Indices) insertCount(url, term string, sentenceID int) error {
	urlID, _, err := mm.getURLID(url)
	if err != nil {
		log.Println("Error getting url id for ", url)
	}
	termID, _, err := mm.getTermID(term)
	if err != nil {
		log.Println("Error getting term id for ", url)
	}

	var hits int
	err = mm.selCount.QueryRow(urlID, termID, sentenceID).Scan(&hits)
	if err == nil {
		// If word exists, update hits table
		hits++ // Incrementing number of hits

		_, err = mm.updHits.Exec(hits, urlID, termID, sentenceID)
		if err != nil {
			log.Fatalf("Error updating hits table: %v", err)
			return err
		}
		// fmt.Printf("Successfully updated hits to %d for URL: %v and Term: %v\n", hits, url, term)
	} else {
		// If word doesn't exist, insert into hits table
		if err == sql.ErrNoRows {
			hits = 1

			_, err = mm.insHits.Exec(urlID, termID, sentenceID)
			if err != nil {
				log.Fatalf("Error inserting into hits table %v", err)
				return err
			}
		}
		// fmt.Printf("Successfully added to hits - Hit: %d, URL: %v, Term: %v\n", hits, url, term)
	}
	return nil
}

// Function to get the word count of the term in the doc
func (mm *Indices) getCount(url, term string) (int, error) {
	var hits int

	urlID, _, err1 := mm.getURLID(url)
	if err1 != nil {
		log.Println("Error getting url id in getCount: ", err1)
		return -1, err1
	}
	termID, _, err2 := mm.getTermID(term)
	if err2 != nil {
		log.Println("Error getting url id in getCount: ", err2)
		return -1, err2
	}
	sentenceID, err3 := mm.getSentenceID(term, url)
	if err3 != nil {
		log.Println("Error getting sentence id in getCount: ", err3)
		return -1, err3
	}

	// Word and url in the hits table
	err := mm.selCount.QueryRow(urlID, termID, sentenceID).Scan(&hits)
	if err == nil {
		return hits, nil
	} else {
		return 0, nil
	}
}

func (mm *Indices) getBigramCount(url, term1, term2 string) (int, error) {
	var hits int

	// Get url and term ids
	urlID, _, err1 := mm.getURLID(url)
	if err1 != nil {
		log.Println("Error getting url id in getBigramCount: ", err1)
		return -1, err1
	}
	firstTermID, _, err2 := mm.getTermID(term1)
	if err2 != nil {
		log.Println("Error getting term id in getBigramCount: ", err2)
		return -1, err2
	}
	secondTermID, _, err3 := mm.getTermID(term2)
	if err3 != nil {
		log.Println("Error getting term id in getBigramCount: ", err3)
		return -1, err3
	}
	// fmt.Println(urlID, firstTermID, secondTermID)
	sentenceID, err4 := mm.getSentenceIDBigram(url, term1, term2)
	if err4 != nil {
		log.Println("Error getting sentence id in getBigramCount: ", err4)
		return -1, err4
	}

	// Get count in the bigram table
	err := mm.selCountBigram.QueryRow(urlID, firstTermID, secondTermID, sentenceID).Scan(&hits)
	if err == nil {
		return hits, nil
	} else {
		return 0, nil
	}
}

// Function to get the total words in a doc
func (mm *Indices) getWordsInDoc(url string) (int, error) {
	urlID, _, err := mm.getURLID(url)
	if err != nil {
		return 0, err
	}
	var total_WID int
	if err := mm.selSumCount.QueryRow(urlID).Scan(&total_WID); err != nil {
		return 0, err
	}
	return total_WID, nil
}

func (mm *Indices) getWordsInDocBigram(url string) (int, error) {
	urlID, _, err := mm.getURLID(url)
	if err != nil {
		return 0, err
	}
	var total_WID int
	if err := mm.selSumCountBigram.QueryRow(urlID).Scan(&total_WID); err != nil {
		return 0, err
	}
	return total_WID, nil
}

// Get total documents with the term in it
func (mm *Indices) getTotalDocWTerm(term string) (int, error) {
	termID, _, err := mm.getTermID(term)
	if err != nil {
		log.Printf("Error getting term id for %v: %v", term, err)
		return 0, err
	}
	var totalWordsInDoc int
	if err := mm.selCountTermID.QueryRow(termID).Scan(&totalWordsInDoc); err == nil {
		return totalWordsInDoc, nil
	} else {
		log.Printf("Error querying to get total words in doc: %v", err)
		return 0, err
	}
}

// Get total documents with the bigram term in it
func (mm *Indices) getTotalDocWTermBigram(first, second string) (int, error) {
	firstTermID, _, err2 := mm.getTermID(first)
	if err2 != nil {
		log.Println("Error getting term id in getTotalDocWTermBigram: ", err2)
		return -1, err2
	}
	secondTermID, _, err3 := mm.getTermID(second)
	if err3 != nil {
		log.Println("Error getting term id in getTotalDocWTermBigram: ", err3)
		return -1, err3
	}

	var totalWordsInDoc int
	if err := mm.selCountTermIDBigram.QueryRow(firstTermID, secondTermID).Scan(&totalWordsInDoc); err == nil {
		return totalWordsInDoc, nil
	} else {
		log.Printf("Error querying to get total words in doc bigram: %v", err)
		return 0, err
	}
}
