Markdown
# URL Shortener

This is a simple URL shortener application written in Go. It allows you to shorten URLs and provides a form for easy use.

## Features

*   Shorten URLs with custom short keys.
*   Redirect to the original URL when accessing the shortened URL.
*   Simple and user-friendly web interface.

## How to Use

1.  **Run the application:**
    *   Clone the repository: `git clone https://github.com/Nathene/url-shortener.git`
    *   Build and run the application: `go build && ./url-shortener`

2.  **Access the URL shortener:**
    *   Modify your hosts file:
        *   Open the hosts file as administrator (located at `C:\Windows\System32\drivers\etc\hosts`).
        *   Add the following line at the end of the file: `127.0.0.1 go`
        *   Save the hosts file.
    *   Open your web browser and go to `go/shorten`

3.  **Shorten a URL:**
    *   Enter the URL you want to shorten in the "Enter a URL" field.
    *   Enter a custom short key in the "Enter a short key" field.
    *   Click the "Shorten" button.

4.  **Access the shortened URL:**
    *   The shortened URL will be displayed in the "Shortened URL" field.
    *   You can now access the original URL by visiting `go/your_short_key`

## Example

**Original URL:** `https://www.example.com/very/long/url/with/many/parameters`

**Short Key:** `shortkey`

**Shortened URL:** `go/shortkey`

## Dockerization

You can also run this application using Docker - which is the preferred option:

1.  **Build the Docker image:**
    ```bash
    docker-compose build
    ```
2.  **Run the Docker container:**
    ```bash
    docker-compose up -d

The examples will work the exact same way, you just wont have to worry about constantly having your server up.