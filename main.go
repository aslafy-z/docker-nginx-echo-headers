package main

import (
        "crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

var hostname string
var randomString string
var delay int
var listenAddr string

func generateRandomString(length int) string {
   b := make([]byte, length)
   _, err := rand.Read(b)
   if err != nil {
      panic(err)
   }
   return base64.StdEncoding.EncodeToString(b)
}

func logMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, randomString)
	time.Sleep(time.Duration(delay * time.Second))
	return
}

func main() {
	host, err := os.Hostname()
	hostname = host
	if err != nil {
		panic(err)
	}
	if os.Getenv("ECHO_BYTES") != "" {
		randomString = generateRandomString(strconv.Atoi(os.Getenv("ECHO_BYTES")))
	}
	if os.Getenv("ECHO_DELAY") != "" {
		delay, _ = strconv.Atoi(os.Getenv("ECHO_DELAY"))
	}
	listenAddr = os.Getenv("ECHO_ADDR")
	if os.Getenv("ECHO_ADDR") == "" {
		listenAddr = ":8080"
	}


	log.Printf("Listening on %s\n", listenAddr)

	http.HandleFunc("/", echoHandler)
	http.ListenAndServe(listenAddr, logMiddleware(http.DefaultServeMux))
}
