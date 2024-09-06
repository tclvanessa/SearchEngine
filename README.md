# Search Engine Web Crawler

This project is a simple search engine web crawler built using Go. It crawls websites, stores data in a relational database, and allows users to search for terms. The search engine implements basic information retrieval techniques like stemming, stop-word omission, and TF-IDF scoring to deliver relevant search results.

## Features

- **Web Crawler**: The web crawler scrapes websites, extracting data and storing it in a relational database.
- **Search Engine**: Users can search for terms through a web interface (styled with CSS), and the most relevant links containing the search terms are displayed.
- **Information Retrieval**: 
  - **Snowball Stemmer**: Reduces words to their root form.
  - **Stop Word Omission**: Common words (e.g., "the", "is") are excluded from the index.
  - **TF-IDF**: Uses Term Frequency-Inverse Document Frequency to rank search results.
- **Snippet Generation**: Provides context for the search results by showing relevant excerpts from the page.

## How It Works

1. **Crawling**: The crawler visits pages and extracts textual content.
2. **Data Storage**: The extracted content is stored in a relational database.
3. **Indexing**: Information is processed using the Snowball stemmer and stop words are removed.
4. **Searching**: Users can input search queries, which are matched against the indexed content using TF-IDF.
5. **Results**: The most relevant pages are shown along with snippets that highlight the search term's occurrence.

## Technologies

- **Backend**: Go
- **Frontend**: HTML, CSS
- **Database**: Relational database (e.g., MySQL, PostgreSQL, or SQLite)
- **Information Retrieval Techniques**: 
  - Snowball Stemmer
  - Stop Word Omission
  - TF-IDF

## Prerequisites

Before running the project, ensure you have:

- Go installed
- A relational database set up (MySQL, PostgreSQL, or SQLite)
- Git for version control

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/search-engine-web-crawler.git
   cd search-engine-web-crawler
   
2. Install dependencies:
  ```bash
  go get ./...
  ```
3. Set up your database:

   - Configure your database connection in the `config.go` file by providing the necessary database credentials, such as the host, port, username, password, and database name.

   Example configuration in `config.go`:

   ```go
   package config

   import (
       "database/sql"
       _ "github.com/go-sql-driver/mysql"
   )

   func ConnectDB() (*sql.DB, error) {
       db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/dbname")
       if err != nil {
           return nil, err
       }
       return db, nil
   }

4. Run the web crawler and search engine:
  To start the crawler and web server:
  ```bash
  go run main.go
  ```
  This will start the server on localhost:8080. You can now access the search interface through your browser.

## Usage

1. Open a web browser and navigate to http://localhost:8080.

2. Use the search bar to enter a query.

3. The search results will display the most relevant links based on your query, including a snippet of the surrounding text.

## Database Configuration: 

  - Update the database connection parameters in the config.go file.
  - Crawler Settings: Modify crawling parameters such as depth and URL patterns in the crawler.go file.
