package main

import (
	"bytes"
	"log"
	"strings"

	"golang.org/x/net/html"

	"gopkg.in/neurosnap/sentences.v1/english"
)

// Function for extracting in dbs
func (mm *Indices) extractDB(inC chan string, chIn chan DownloadResult, chOut chan ExtractResult) {
	// fmt.Println("extracting!")

	// Range over input channel
	for dl := range chIn {
		// words := make([]string, 0)
		sentences := make([]string, 0)
		var title string

		doc, err := html.Parse(bytes.NewReader(dl.body))
		if err != nil {
			log.Fatal("Extract Parse error: ", err)
		}

		// Initialize sentence tokenizer
		tokenizer, err := english.NewSentenceTokenizer(nil)
		if err != nil {
			log.Println("Sentence Tokenizer Error: ", err)
		}

		var f func(*html.Node)
		f = func(n *html.Node) {
			switch n.Type {
			case html.ElementNode:
				if n.Data == "title" {
					title = n.FirstChild.Data
				}
			case html.TextNode:
				p := n.Parent
				if p.Type == html.ElementNode && (p.Data != "style" && p.Data != "script") {
					sentence := tokenizer.Tokenize(strings.TrimSpace(n.Data))

					for _, s := range sentence {
						sentences = append(sentences, s.Text)
					}

					// newWords := strings.FieldsFunc(n.Data, func(r rune) bool {
					// 	return !unicode.IsLetter(r) && !unicode.IsNumber(r)
					// })
					// for _, w := range newWords {
					// 	stem, err := snowball.Stem(w, "english", true)
					// 	if err != nil {
					// 		log.Println("Failed to stem", w, err)
					// 	} else {
					// 		// Check if the stemmed word is not a stop word
					// 		if _, exists := mm.stopWordsMap[stem]; !exists {
					// 			words = append(words, stem)
					// 		}
					// 	}
					// }
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)

		chOut <- ExtractResult{dl.url /*words, */, title, sentences}
	}
}
