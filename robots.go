package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (mm *Indices) robotsMap(base string) {
	mm.RobotRecords = make(map[string]Rules)

	// Make robots.txt url from base url with helper function
	robotsUrl := makeRobotsURL(base)

	// Get robots.txt url
	if rsp, err1 := http.Get(robotsUrl); err1 == nil {
		// Read body of robots.txt url
		if robotData, err2 := io.ReadAll(rsp.Body); err2 == nil {
			// Split body of robots into lines
			lines := strings.Split(string(robotData), "\n")
			user := ""

			// Range over each line
			for _, line := range lines {
				// If on User-agent line,
				if strings.HasPrefix(line, "User-agent:") {
					tUser := strings.TrimSpace(strings.TrimPrefix(line, "User-agent:"))
					user = strings.ReplaceAll(tUser, "*", ".*")
				}

				// If on Allow, Disallow, or delay line
				if strings.HasPrefix(line, "Allow:") {
					allowRules := mm.RobotRecords[user]
					aPath := strings.TrimSpace(strings.TrimPrefix(line, "Allow:"))
					allowPath := strings.ReplaceAll(aPath, "*", ".*")
					allowRules.Allow = append(allowRules.Allow, allowPath)
					mm.RobotRecords[user] = allowRules
				} else if strings.HasPrefix(line, "Disallow:") {
					disallowRules := mm.RobotRecords[user]
					dPath := strings.TrimSpace(strings.TrimPrefix(line, "Disallow:"))
					disallowPath := strings.ReplaceAll(dPath, "*", ".*")
					disallowRules.Disallow = append(disallowRules.Disallow, disallowPath)
					mm.RobotRecords[user] = disallowRules
				} else if strings.HasPrefix(line, "crawl-delay:") {
					delayNum := strings.TrimSpace(strings.TrimPrefix(line, "crawl-delay:"))
					if num, err := strconv.Atoi(delayNum); err == nil {
						delayRules := mm.RobotRecords[user]
						delayRules.delay = num
						mm.RobotRecords[user] = delayRules
					}
				}

				// If on sitemap line
				if strings.HasPrefix(line, "Sitemap:") {
					smRule := mm.RobotRecords[user]
					smURL := strings.TrimSpace(strings.TrimPrefix(line, "Sitemap:"))
					smRule.Sitemap = smURL
					mm.RobotRecords[user] = smRule
				}
			}
		} else {
			log.Println("Error in reading body: ", err2)
		}
	} else {
		log.Println("Error in getting url: ", err1)
	}
}

func makeRobotsURL(base string) string {
	// Parse base url
	pURL, err := url.Parse(base)
	if err != nil {
		log.Fatal("Error parsing base url for robots.txt: ", err)
	}
	// Make robots.txt url out of base url
	robotsUrl := pURL.Scheme + "://" + pURL.Host + "/robots.txt"

	return robotsUrl
}

// Function to get the user agent
func (mm *Indices) getUserAgent(base string) string {
	// Make robots.txt url from base url with helper function
	robotsUrl := makeRobotsURL(base)

	// Get robots.txt url
	rsp, err1 := http.Get(robotsUrl)
	if err1 != nil {
		log.Fatal("Error getting robots.txt url: ", err1)
	}

	// Read body of robots.txt url
	robotData, err2 := io.ReadAll(rsp.Body)
	if err2 != nil {
		log.Fatal("Error reading body of robots.txt url: ", err1)
	}

	// Split body of robots into lines
	lines := strings.Split(string(robotData), "\n")
	user := ""

	// Range over each line
	for _, line := range lines {
		// If on User-agent line,
		if strings.HasPrefix(line, "User-agent:") {
			tUser := strings.TrimSpace(strings.TrimPrefix(line, "User-agent:"))
			user = strings.ReplaceAll(tUser, "*", ".*")
			return user
		}
	}
	return ""
}
