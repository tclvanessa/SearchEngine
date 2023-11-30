package main

import (
	"testing"
)

// Bigram Test Case
func TestCase1(t *testing.T) {
	tests := []struct {
		term string
		want Hits
	}{
		{"report authors",
			Hits{
				{"https://openai.com/research/confidence-building-measures-for-artificial-intelligence", 0.43402298850574716},
				{"https://openai.com/research/gpts-are-gpts", 0.4104347826086957},
				{"https://openai.com/research/frontier-ai-regulation", 0.3814141414141414},
				{"https://openai.com/research/improving-verifiability", 0.2448767833981842},
				{"https://openai.com/research/forecasting-misuse", 0.12356020942408379},
			},
		},
	}

	// Open connection to the original database
	// db, err := sql.Open("sqlite3", "sqlindex.db")
	// if err != nil {
	// 	log.Fatal("Error opening database: ", err)
	// }
	// defer db.Close()

	// Create an instance of Indices with the existing database connection
	// mmDB := &Indices{db: db}

	// mmDB.totalDoc = 472

	mmDB := NewIndices()
	// url := "https://openai.com/"

	// serveDB(mmDB)
	// mmDB.robotsMap(url)
	// mmDB.crawlDB(url, 40, 150, 200)

	for _, tc := range tests {
		tfidf, _ := mmDB.BigramTfIdfDB(tc.term)
		for i := range tfidf {
			if tfidf[i].URL != tc.want[i].URL || tfidf[i].Score != tc.want[i].Score {
				t.Errorf("For %v:\nGot: %v\nExpected: %v\n", tc.term, tfidf, tc.want)
				break // Exit the loop early if a difference is found
			}
		}
	}
}

// Wildcard Test Case
func TestCase2(t *testing.T) {
	tests := []struct {
		term string
		want Hits
	}{
		{"wood",
			Hits{
				{"https://openai.com/research/vpt", 0.9662231320368475},
				{"https://openai.com/blog/chatgpt-can-now-see-hear-and-speak", 0.35276532137518685},
				{"https://openai.com/dall-e-3", 0.301404853128991},
				{"https://openai.com/blog/dall-e-2-extending-creativity", 0.295369211514393},
				{"https://openai.com/research/fine-tuning-gpt-2", 0.14717804801995632},
				{"https://openai.com/research/better-language-models", 0.11017740429505134},
			},
		},
	}

	// db, err := sql.Open("sqlite3", "sqlindex.db")
	// if err != nil {
	// 	log.Fatal("Error opening database: ", err)
	// }
	// defer db.Close()

	// Create an instance of Indices with the existing database connection
	// mmDB := &Indices{db: db}

	// mmDB.totalDoc = 472

	mmDB := NewIndices()
	// url := "https://openai.com/"

	// // serveDB(mmDB)
	// mmDB.robotsMap(url)
	// mmDB.crawlDB(url, 40, 150, 200)

	for _, tc := range tests {
		tfidf, _ := mmDB.WildcardTfIdfDB(tc.term)
		for i := range tfidf {
			if tfidf[i].URL != tc.want[i].URL || tfidf[i].Score != tc.want[i].Score {
				t.Errorf("For %v:\nGot: %v\nExpected: %v\n", tc.term, tfidf, tc.want)
				break // Exit the loop early if a difference is found
			}
		}
	}
}
