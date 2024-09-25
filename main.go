package main

import (
	"fmt"
	"net/http"
)

// URLShortener will store a map of short urls as keys, with the longer urls as the values.
type URLShortener struct {
	urls map[string]string
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

	us.urls[customShortKey] = originalURL

	//shortenedURL := fmt.Sprintf(baseURL+"%s", customShortKey)

	handleGET(us, w, r)
}

// handleGET will show the basic UI for shortening a url
func handleGET(us *URLShortener, w http.ResponseWriter, r *http.Request) {
	// Display the HTML form
	w.Header().Set("Content-Type", "text/html")
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

	for k, v := range us.urls {
		responseHTML += fmt.Sprintf(`
                <tr>
                    <td>%s</td>
                    <td><a href="/%s">%s</a></td> 
                </tr>
        `, v, k, "go/"+k) // Note: v is the long URL, k is the short key
	}

	responseHTML += `
            </tbody>
        </table>
    `
	fmt.Fprintf(w, responseHTML)
}

// HandRedirection will use the shortened url and get the value from the map inside URLShortener
func (us *URLShortener) HandleRedirection(w http.ResponseWriter, r *http.Request) {
	// ignores the intitial '/'
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

	http.HandleFunc("/shorten", shortener.HandleShorten) // Handle shortening at /shorten
	http.HandleFunc("/", shortener.HandleRedirection)

	http.ListenAndServe(port, nil)
}
