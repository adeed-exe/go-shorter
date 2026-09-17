package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	registerPage = `
    <!DOCTYPE html>
    <html>
        <head><title>url shortener</title></head>
        <body style = "display: flex; justify-content: center; align-items: center; height: 100vh;">
            <div style = "text-align: center;">
                <h2>url shortener</h2>
                <form method = "POST" action = "/register">
                    <input type = "url" name = "url" placeholder = "enter your url" required>
                    <button type = "submit">shorten</button>
                </form>
            </div>
        </body>
    </html>`

	outputPage = `
    <!DOCTYPE html>
    <html>
        <head><title>short link output</title></head>
        <body style = "display: flex; justify-content: center; align-items: center; height: 100vh;">
            <div style = "text-align: center;">
                <h3>your short link</h3>
                <p>%s</p>
                <p><a href="/register">shorten another url</a></p>
            </div>
        </body>
    </html>`
)

var (
	count   uint64 = 0
	urlToID        = make(map[string]string)
	idToURL        = make(map[string]string)
)

func generateID() string {
	current := count
	count++

	var id string
	for range 7 {
		r := current % 62
		id = string(charset[r]) + id
		current /= 62
	}
	return id
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		url := r.FormValue("url")

		id, found := urlToID[url]
		if !found {
			id = generateID()
			urlToID[url] = id
			idToURL[id] = url
		}

		http.Redirect(w, r, "/output?id="+id, http.StatusSeeOther)
		return
	}

	fmt.Fprint(w, registerPage)
}

func outputHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	fmt.Fprintf(w, outputPage, "http://"+r.Host+"/"+id)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")

	if id == "" {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	url, found := idToURL[id]
	if found {
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	http.Error(w, "Short link not registered", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/output", outputHandler)
	http.HandleFunc("/", redirectHandler)

	port := process.env.PORT
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on port %s", port)
	http.ListenAndServe(":"+port, nil)
}
