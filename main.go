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

	env "github.com/caarlos0/env/v8"
)

var (
	
        echoContext       bool
	randomBytesCount  int
	randomBytes       string
	echoDelay         time.Duration
	listenAddr        string
)

type config struct {
	ShowContext       bool `env:"ECHO_CONTEXT" envDefault:"false"`
	RandomBytes       int `env:"ECHO_BYTES" envDefault:"0"`
	Delay             time.Duration `env:"ECHO_DELAY" envDefault:"0s"`
	ListenAddress     string `env:"LISTEN_ADDRESS,required" envDefault:":8080"`
	Debug             bool   `env:"DEBUG" envDefault:"false"`
}

var (
	cfg         config
	hostname    string
	randomBytes string
)

func randomBytes(size int) (string, error) {
   b := make([]byte, size)
   _, err := rand.Read(b)
   if err != nil {
      return nil, err
   }
   return base64.StdEncoding.EncodeToString(b), nil
}

func logMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(cfg.Delay))
	w.Header().Set("Content-Type", "text/plain")
	if echoContext {
		fmt.Fprintln(w, "X-Echo-Date:", time.Now().String())
	        fmt.Fprintf(w, "X-Echo-Delay: %s\n", cfg.Delay)
		fmt.Fprintln(w, "X-Echo-Hostname:", hostname)
		fmt.Fprintln(w, "X-Echo-Method:", r.Method)
		fmt.Fprintln(w, "X-Echo-Proto:", r.Proto)
	        fmt.Fprintf(w, "X-Echo-RandomBytes: %d\n", cfg.RandomBytes)
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
	fmt.Fprintln(w, randomBytes)
	return
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

	// cache values
	hostname, _ = os.Hostname()
	if cfg.RandomBytes > 0 {
		randomString, err = randomBytes(cfg.RandomBytes)
		if err != nil {
			log.Fatalf("error: generate random bytes: %+v\n", err)
		}
	}

	log.Printf("Listening on %s\n", cfg.ListenAddr)
	http.HandleFunc("/-/ready", handleReadinessRequest)
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(cfg.ListenAddress, logMiddleware(http.DefaultServeMux))
}
