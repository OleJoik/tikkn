package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, `<!DOCTYPE html>
<html>
  <head><title>Hello</title></head>
  <body><h1>Hello from Go!</h1></body>
</html>`)
}

func main() {
	httpPort := getEnv("HTTP_PORT", "8080")
	httpsPort := getEnv("HTTPS_PORT", "8443")

	go func() {
		redirect := func(w http.ResponseWriter, r *http.Request) {
			host, _, err := net.SplitHostPort(r.Host)
			if err != nil {
				host = r.Host
			}

			target := "https://" + host

			if httpsPort != "443" {
				target += ":" + httpsPort
			}

			target += r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		}

		addr := ":" + httpPort
		fmt.Printf("Redirecting HTTP on %s → HTTPS on %s\n", addr, httpsPort)

		if err := http.ListenAndServe(addr, http.HandlerFunc(redirect)); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP redirect server failed: %v\n", err)
			os.Exit(1)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	addr := ":" + httpsPort
	fmt.Printf("Serving HTTPS on %s\n", addr)
	if err := http.ListenAndServeTLS(addr,
		mustGetEnv("CERT_PATH"),
		mustGetEnv("KEY_PATH"),
		mux); err != nil {
		panic(err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	value, found := os.LookupEnv(key)
	if !found {
		panic(fmt.Sprintf("%s is a required env var!", key))
	}
	return value
}
