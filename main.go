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

	"github.com/ahuigo/gofnext"
	env "github.com/caarlos0/env/v8"
)

type config struct {
	ShowContext       bool          `env:"ECHO_CONTEXT" envDefault:"false"`
	RandomBytesLength int           `env:"ECHO_RAND_BYTES" envDefault:"0"`
	Delay             time.Duration `env:"ECHO_DELAY" envDefault:"0s"`
	ListenAddress     string        `env:"ECHO_ADDR,required" envDefault:":8080"`
}

var (
	cfg                  config
	hostname             string
	getRandomBytesCached func(int) string
)

func getRandomBytes(size int) string {
	if size < 1 {
		return ""
	}
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func logMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	delay := cfg.Delay
	if val := r.URL.Query().Get("delay"); val != "" {
		duration, err := time.ParseDuration(val)
		if err != nil {
			http.Error(w, "Invalid delay value", http.StatusBadRequest)
			return
		}
		delay = duration
	}

	randomBytesLength := cfg.RandomBytesLength
	if val := r.URL.Query().Get("bytes"); val != "" {
		length, err := strconv.Atoi(val)
		if err != nil {
			http.Error(w, "Invalid bytes value", http.StatusBadRequest)
			return
		}
		randomBytesLength = length
	}
	randomBytes := getRandomBytesCached(randomBytesLength)

	time.Sleep(time.Duration(delay))
	w.Header().Set("Content-Type", "text/plain")
	if cfg.ShowContext {
		fmt.Fprintln(w, "X-Echo-Date:", time.Now().String())
		fmt.Fprintf(w, "X-Echo-EffectiveDelay: %s\n", delay)
		fmt.Fprintln(w, "X-Echo-Hostname:", hostname)
		fmt.Fprintln(w, "X-Echo-Method:", r.Method)
		fmt.Fprintln(w, "X-Echo-Proto:", r.Proto)
		fmt.Fprintf(w, "X-Echo-EffectiveRandomBytesLength: %d\n", randomBytesLength)
		fmt.Fprintln(w, "X-Echo-RemoteAddr:", r.RemoteAddr)
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
	if randomBytes != "" {
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, randomBytes)
	}
}

// handleReadinessRequest handles incoming readiness requests
func handleReadinessRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	// retrieve configuration
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error: configuration parsing: %+v\n", err)
	}

	getRandomBytesCached = gofnext.CacheFn1(getRandomBytes)
	hostname, _ = os.Hostname()

	log.Printf("Listening on %s\n", cfg.ListenAddress)
	http.HandleFunc("/-/ready", handleReadinessRequest)
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(cfg.ListenAddress, logMiddleware(http.DefaultServeMux))
}
