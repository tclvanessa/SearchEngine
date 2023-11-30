package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
)

type SitemapIndex struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []string `xml:"url>loc"`
}

func (mm *Indices) sitemapURLs(base string) []string {
	// Make robots.txt url from base url with helper function
	robotURL := makeRobotsURL(base)

	userAgent := mm.getUserAgent(robotURL)
	sitemapURL := mm.RobotRecords[userAgent].Sitemap

	// Get sitemap url
	rsp, err1 := http.Get(sitemapURL)
	if err1 != nil {
		log.Fatal("Error getting sitemap url: ", err1)
	}
	defer rsp.Body.Close()

	// Read the response body
	body, err2 := io.ReadAll(rsp.Body)
	if err2 != nil {
		log.Fatal("Error reading body of sitemap url: ", err2)
	}

	// Create a variable to store the unmarshaled data
	var sitemap SitemapIndex

	// Unmarshal the XML data into the data structure
	err3 := xml.Unmarshal(body, &sitemap)
	if err3 != nil {
		log.Fatal("Error unmarshalling xml data: ", err3)
	}
	return sitemap.URLs
}
