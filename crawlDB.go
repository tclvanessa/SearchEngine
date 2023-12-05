package main

import (
	"fmt"
	"time"
)

// var wg sync.WaitGroup

// Crawl function using concurrency
func (mm *Indices) crawlDB(seed string, downloaders int, extractors int, indexers int) {
	// Size of input and output channels
	const dlInCSize = 800
	const dlOutCSize = 850
	const exOutCSize = 550

	// Make channels for input, output, extract output, and quit channel
	dlInC := make(chan string, dlInCSize)
	defer close(dlInC)

	dlOutC := make(chan DownloadResult, dlOutCSize)
	defer close(dlOutC)

	exOutC := make(chan ExtractResult, exOutCSize)
	defer close(exOutC)

	quitC := make(chan struct{}, 1)

	// Increment the WaitGroup for each indexer goroutine
	// wg.Add(indexers)

	// Get the list of URLs from the sitemap
	sitemapURLs := mm.sitemapURLs(seed)

	// Start the crawl with URLs from the sitemap
	for _, url := range sitemapURLs {
		// fmt.Println(url)
		dlInC <- url // Feed each URL into the download input channel
	}

	// Goroutine to run concurrently with crawl to monitor if index stops growing
	go func() {
		prevDoc := 0 // Keep track of documents

		// Go until quit channel
		for {
			time.Sleep(7 * time.Second) // Sleeping for a second to wait for index to grow
			outstanding := len(dlInC) + len(dlOutC) + len(exOutC)

			// If the documents are the same as last time (means there's nothing else has been crawled),
			// and all of the channels are empty, then quit
			// Some sites have lots of docs, so stop crawling after 5,000 docs
			fmt.Println(mm.totalDoc)
			if (prevDoc == mm.totalDoc && outstanding == 0) || mm.totalDoc >= 5000 {
				quitC <- struct{}{}
				break
			} else {
				prevDoc = mm.totalDoc
			}
		}
	}()

	// Multiple workers to make crawl faster
	for i := 0; i < downloaders; i++ {
		// Take in input channel and output body/url
		go mm.downloadDB(dlInC, dlOutC)
	}
	for i := 0; i < extractors; i++ {
		// Take in body/url from download to extract words urls
		// Also send clean urls not crawled yet into download queue
		go mm.extractDB(dlInC, dlOutC, exOutC)
	}
	for i := 0; i < indexers; i++ {
		// defer wg.Done()
		go mm.indexDB(exOutC) // Send in extract channel to index
	}
	// Quit channel
	<-quitC

	fmt.Println("Total docs: ", mm.totalDoc)
}
