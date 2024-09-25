package main

import (
	"fmt"
	"net/http"
)

type URLShortener struct {
	urls map[string]string
}

const (
	baseURL = "http://localhost:80/"
	port    = ":80"
)

func (us *URLShortener) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handleGET(w, r)
		return
	}
	if r.Method == "POST" {
		handlePOST(us, w, r)
		return
	}
	http.Error(w, "method not supported: "+r.Method, http.StatusBadRequest)
}

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

	us.urls[customShortKey] = originalURL

	shortenedURL := fmt.Sprintf(baseURL+"%s", customShortKey)

	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
        <h2>URL Shortener</h2>
        <p>Original URL: <input type="text" value="%s" readonly></p>
        <p>Shortened URL: <input type="text" value="%s" readonly></p>
        <form method="post" action="/shorten"> 
            <input type="text" name="url" placeholder="Enter a URL">
            <input type="text" name="shortkey" placeholder="Enter a short key"> 
            <input type="submit" value="Shorten">
        </form>
    `, originalURL, shortenedURL)
	fmt.Fprintf(w, responseHTML)
}

func handleGET(w http.ResponseWriter, r *http.Request) {
	// Display the HTML form
	w.Header().Set("Content-Type", "text/html")
	responseHTML := `
		<h2>URL Shortener</h2>
		<form method="post" action="/shorten"> 
			<input type="text" name="url" placeholder="Enter a URL">
			<input type="text" name="shortkey" placeholder="Enter a short key"> 
			<input type="submit" value="Shorten">
		</form>
	`
	fmt.Fprintf(w, responseHTML)
}

func (us *URLShortener) HandleRedirection(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[1:]
	fmt.Println(us.urls)
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	originalURL, found := us.urls[shortKey]
	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func main() {
	shortener := &URLShortener{
		urls: make(map[string]string),
	}

	http.HandleFunc("/shorten", shortener.HandleShorten) // Handle shortening at root path
	http.HandleFunc("/", shortener.HandleRedirection)

	fmt.Println("URL shortener is running on " + port)
	fmt.Println(baseURL + "shorten")
	http.ListenAndServe(port, nil)
}
