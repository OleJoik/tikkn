package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/OleJoik/tikkn/middleware"
)

func main() {
	loadEnv(".env")

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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "./public" + r.URL.Path
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			http.ServeFile(w, r, path)
			return
		}

		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			http.ServeFile(w, r, "./public/index.html")
			return
		}

		http.NotFound(w, r)
	})

	handler := middleware.Logging(mux)

	addr := ":" + httpsPort
	fmt.Printf("Serving HTTPS on %s\n", addr)
	if err := http.ListenAndServeTLS(addr,
		getEnv("RUNTIME_CERT_PATH", "/app/certs/cert.pem"),
		getEnv("RUNTIME_KEY_PATH", "/app/certs/privkey.pem"),
		handler); err != nil {
		panic(err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func loadEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines or comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		// Strip optional quotes
		val = strings.Trim(val, `"'`)

		os.Setenv(key, val)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
