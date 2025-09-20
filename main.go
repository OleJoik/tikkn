package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Set content type so the browser knows it’s HTML
	w.Header().Set("Content-Type", "text/html")

	// Send a minimal HTML page
	fmt.Fprintln(w, `<!DOCTYPE html>
<html>
  <head><title>Hello</title></head>
  <body><h1>Hello from Go!</h1></body>
</html>`)
}

func main() {
	http.HandleFunc("/", handler)

	// Listen on port 8080
	fmt.Println("Serving on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
