package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// URLShortener will store a map of short urls as keys, with the longer urls as the values.
type URLShortener struct {
	urls map[string]string
	db   *sql.DB
}

const (
	baseURL = "http://localhost:80/"
	port    = ":80"
)

func (us *URLShortener) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handleGET(us, w, r)
		return
	}
	if r.Method == "POST" {
		handlePOST(us, w, r)
		return
	}
	http.Error(w, "method not supported: "+r.Method, http.StatusBadRequest)
}

// handlePOST will do all checks on the urls, make sure there is no duplicates and then send the complete form.
func handlePOST(us *URLShortener, w http.ResponseWriter, r *http.Request) {
	originalURL := r.FormValue("url")

	if originalURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	customShortKey := r.FormValue("shortkey")
	if customShortKey == "" {
		http.Error(w, "Shortkey parameter is missing", http.StatusBadRequest)
		return
	}

	_, exists := us.urls[customShortKey]
	if exists {
		http.Error(w, fmt.Sprintf("shortkey %s already exists.", customShortKey), http.StatusBadRequest)
		return
	}

	_, err := us.db.Exec("INSERT INTO urls (short_key, original_url) VALUES ($1, $2)", customShortKey, originalURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting URL: %v", err), http.StatusInternalServerError)
		return
	}

	us.urls[customShortKey] = originalURL

	//shortenedURL := fmt.Sprintf(baseURL+"%s", customShortKey)

	handleGET(us, w, r)
}

// handleGET will show the basic UI for shortening a url
func handleGET(us *URLShortener, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// 1. Query the database to get all URLs
	rows, err := us.db.Query("SELECT short_key, original_url FROM urls")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying URLs: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// 2. Create a slice to hold the URL data
	var urls []struct {
		ShortKey    string `json:"short_key"`
		OriginalURL string `json:"original_url"`
	}

	// 3. Iterate through the rows and populate the slice
	for rows.Next() {
		var url struct {
			ShortKey    string `json:"short_key"`
			OriginalURL string `json:"original_url"`
		}
		if err := rows.Scan(&url.ShortKey, &url.OriginalURL); err != nil {
			http.Error(w, fmt.Sprintf("Error scanning URL row: %v", err), http.StatusInternalServerError)
			return
		}
		urls = append(urls, url)
	}

	// 4. Generate the HTML with the URL data
	responseHTML := `
        <h2>URL Shortener</h2>
        <form method="post" action="/shorten">
            <input type="text" name="url" placeholder="Enter a URL">
            <input type="text" name="shortkey" placeholder="go/">
            <input type="submit" value="Shorten">
        </form>
        <table>
            <thead>
                <tr>
                    <th>Long URL</th>
                    <th>Short URL</th>
                </tr>
            </thead>
            <tbody>
    `

	for _, item := range urls {
		responseHTML += fmt.Sprintf(`
                <tr>
                    <td>%s</td>
                    <td><a href="/%s">go/%s</a></td>
                </tr>
        `, item.OriginalURL, item.ShortKey, item.ShortKey)
	}

	responseHTML += `
            </tbody>
        </table>
    `
	fmt.Fprintf(w, responseHTML)
}
func (us *URLShortener) HandleGetURLs(w http.ResponseWriter, r *http.Request) {
	rows, err := us.db.Query("SELECT short_key, original_url FROM urls")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying URLs: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var urls []struct {
		ShortKey    string `json:"short_key"`
		OriginalURL string `json:"original_url"`
	}
	for rows.Next() {
		var url struct {
			ShortKey    string `json:"short_key"`
			OriginalURL string `json:"original_url"`
		}
		if err := rows.Scan(&url.ShortKey, &url.OriginalURL); err != nil {
			http.Error(w, fmt.Sprintf("Error scanning URL row: %v", err), http.StatusInternalServerError)
			return
		}
		urls = append(urls, url)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(urls)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(urls)
}

// HandRedirection will use the shortened url and get the value from the map inside URLShortener
func (us *URLShortener) HandleRedirection(w http.ResponseWriter, r *http.Request) {
	// ignores the intitial '/'
	shortKey := r.URL.Path[1:]

	// Check the in-memory map first
	originalURL, found := us.urls[shortKey]
	if found {
		http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
		return
	}

	// If not found in the map, query the database
	err := us.db.QueryRow("SELECT original_url FROM urls WHERE short_key = $1", shortKey).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Shortened key not found", http.StatusNotFound)
		} else {
			http.Error(w,
				fmt.Sprintf("Error retrieving URL: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Update the map with the retrieved URL for caching
	us.urls[shortKey] = originalURL

	fmt.Println(us.urls)
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	originalURL, found = us.urls[shortKey]
	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func main() {
	db, err := sql.Open("postgres", PG_CONN_STRING)
	if err != nil {
		log.Fatal(err)
	}

	shortener := &URLShortener{
		urls: make(map[string]string),
		db:   db,
	}

	http.HandleFunc("/urls", shortener.HandleGetURLs)
	http.HandleFunc("/add", shortener.HandleShorten)
	http.HandleFunc("/", shortener.HandleRedirection)

	http.ListenAndServe(port, nil)
}
