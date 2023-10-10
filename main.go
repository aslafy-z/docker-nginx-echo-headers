package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

var hostname string
var echoContext bool
var listenAddr string

func logMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	if echoContext {
		fmt.Fprintln(w, "X-Echo-Hostname:", hostname)
		fmt.Fprintln(w, "X-Echo-RemoteAddr:", r.RemoteAddr)
		fmt.Fprintln(w, "X-Echo-Date:", time.Now().String())
		fmt.Fprintln(w, "X-Echo-Proto:", r.Proto)
		fmt.Fprintln(w, "X-Echo-Method:", r.Method)
		fmt.Fprintln(w, "X-Echo-URL:", r.URL)
	}
	var headers []string
	for k, vs := range r.Header {
		for _, v := range vs {
			headers = append(headers, fmt.Sprintf("%s: %s", k, v))
		}
	}
	sort.Strings(headers)
	for _, h := range headers {
		fmt.Fprintln(w, h)
	}
	return
}

func main() {
	host, err := os.Hostname()
	hostname = host
	if err != nil {
		panic(err)
	}
	echoContext = true
	if os.Getenv("ECHO_CONTEXT") == "false" {
		echoContext = false
	}
	listenAddr = os.Getenv("ECHO_ADDR")
	if os.Getenv("ECHO_ADDR") == "" {
		listenAddr = ":8080"
	}

	log.Printf("Listening on %s\n", listenAddr)

	http.HandleFunc("/", echoHandler)
	http.ListenAndServe(listenAddr, logMiddleware(http.DefaultServeMux))
}
