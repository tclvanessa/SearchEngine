package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Freq map[string]int
type Rules struct {
	Allow    []string
	Disallow []string
	delay    int
	Sitemap  string
}
type Indices struct {
	db                     *sql.DB
	insTermStmt            *sql.Stmt
	selTerm                *sql.Stmt
	selTermIDStmt          *sql.Stmt
	selTermIDLikeStmt      *sql.Stmt
	selExisTermStmt        *sql.Stmt
	insURLStmt             *sql.Stmt
	selURLIDStmt           *sql.Stmt
	selURLStmt             *sql.Stmt
	selExisURLStmt         *sql.Stmt
	selURLIDFromTermIDStmt *sql.Stmt
	selURLIDFromBigramStmt *sql.Stmt
	selCount               *sql.Stmt
	selCountTermID         *sql.Stmt
	selCountTermIDBigram   *sql.Stmt
	updHits                *sql.Stmt
	insHits                *sql.Stmt
	selSumCount            *sql.Stmt
	selSumCountBigram      *sql.Stmt
	insBigram              *sql.Stmt
	updBigram              *sql.Stmt
	selCountBigram         *sql.Stmt
	updTitle               *sql.Stmt
	selTitle               *sql.Stmt
	insSentenceStmt        *sql.Stmt
	selSentence            *sql.Stmt
	selSentenceBigram      *sql.Stmt
	selSentenceID          *sql.Stmt
	selBigramSentence      *sql.Stmt
	selBigramSentenceID    *sql.Stmt

	totalDoc     int
	stopWordsMap map[string]struct{}
	RobotRecords map[string]Rules
}

// Factory of Indices
func NewIndices() *Indices {
	res := &Indices{}
	res.totalDoc = 0
	res.stopWordsMap = mapStopWords("stopwords-en.json")
	res.createDB()
	return res
}

type DownloadResult struct {
	body []byte
	url  string
}

type ExtractResult struct {
	url string
	// words     []string
	title     string
	sentences []string
}

// Function to read the JSON stop words file and put it into a map
func mapStopWords(filePath string) map[string]struct{} {
	// Read the JSON file
	jsonBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Create a map to hold the stop words
	stopWordsMap := make(map[string]struct{})

	// Unmarshal the JSON array into the map
	var stopWordsArray []string
	if err := json.Unmarshal(jsonBytes, &stopWordsArray); err != nil {
		log.Fatal(err)
	}

	// Populate the map with stop words
	for _, s := range stopWordsArray {
		stopWordsMap[s] = struct{}{}
	}
	return stopWordsMap
}

func main() {
	fmt.Println("Project 06")

	mm := NewIndices()
	url := "https://openai.com/"

	// Initialize workers for crawl
	const downloaders = 40
	const extractors = 150
	const indexers = 200

	serveDB(mm)
	mm.robotsMap(url)
	mm.crawlDB(url, downloaders, extractors, indexers)

	for {
		time.Sleep(10 * time.Millisecond)
	}
}
