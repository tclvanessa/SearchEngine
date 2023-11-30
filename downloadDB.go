package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

// Function for downloading in databases
func (mm *Indices) downloadDB(inC chan string, outC chan DownloadResult) {
	for url := range inC {
		// fmt.Println("downloadDB working for", url)
		var isAllowed bool
		for _, disallow := range mm.RobotRecords[".*"].Disallow {
			if match, err := regexp.MatchString(disallow, url); err == nil {
				if match {
					fmt.Println("isAllowed made to false in disallow")
					isAllowed = false
					break
				} else {
					isAllowed = true
				}
			}
		}

		// If url is allowed
		if isAllowed {
			// Download contents from url
			rsp, err1 := http.Get(url)
			if err1 != nil {
				log.Fatal("Error downloading url: ", err1)
			}

			// Continue if response status code is 200
			if rsp.StatusCode != 200 {
				log.Fatal("StatusCode not 200")
			}

			// Reads response body
			body, err2 := io.ReadAll(rsp.Body)
			if err2 != nil {
				log.Fatal("Error reading body: ", err2)
			}

			// Send url to output channel if it's not in db
			if ok, err3 := mm.insertURL(url); err3 == nil {
				if !ok {
					outC <- DownloadResult{
						body,
						url,
					}
				}
			}
			// Delay crawl for the amount of time in robots.txt
			// If no crawl-delay is found, use 100ms as default
			if mm.RobotRecords[".*"].delay != 0 {
				delayCrawl := time.Duration(mm.RobotRecords[".*"].delay)
				time.Sleep(delayCrawl * time.Second)
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
