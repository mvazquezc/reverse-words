package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

// Input word
type Input struct {
	Word string `json:"word,omitempty"`
}

// Output reverse word
type Output struct {
	ReverseWord string `json:"reverse_word,omitempty"`
}

var (
	totalWordsReversed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "reversewords_reversed_words_total",
			Help: "Total number of reversed words",
		},
	)
)

var (
	endpointsAccessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "reversewords_endpoints_accessed_total",
			Help: "Total number of accessed to a given endpoint",
		},
		[]string{"accessed_endpoint"},
	)
)

const version = "v0.0.27"

// ReturnRelease returns the release configured by the user
func ReturnRelease(w http.ResponseWriter, r *http.Request) {
	release := getEnv("RELEASE", "NotSet")
	releaseString := "Reverse Words Release: " + release + ". App version: " + version
	_, err := w.Write([]byte(releaseString + "\n"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	endpointsAccessed.WithLabelValues("release").Inc()
}

// ReturnHostname returns the hostname for the node where the app is running
func ReturnHostname(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}
	hostnameString := "Hostname: " + hostname
	_, err = w.Write([]byte(hostnameString + "\n"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	endpointsAccessed.WithLabelValues("hostname").Inc()
}

// ReturnHealth returns healthy string, can be used for monitoring pourposes
func ReturnHealth(w http.ResponseWriter, r *http.Request) {
	health := "Healthy"
	_, err := w.Write([]byte(health + "\n"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	endpointsAccessed.WithLabelValues("health").Inc()
}

// ReverseWord returns a reversed word based on an input word
func ReverseWord(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var word Input
	var reverseWord string
	err := decoder.Decode(&word)
	if err != nil {
		// Error EOF means no json data has been sent
		if err.Error() != "EOF" {
			panic(err)
		}
	}
	if len(word.Word) < 1 {
		log.Println("No word detected, sending default reverse word")
		reverseWord = "detceted drow oN"
	} else {
		log.Println("Detected word", word.Word)
		reverseWord = reverse(word.Word)
	}
	log.Println("Reverse word", reverseWord)
	totalWordsReversed.Inc()
	output := Output{reverseWord}
	js, err := json.Marshal(output)
	js = append(js, "\n"...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(js))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	endpointsAccessed.WithLabelValues("reverseword").Inc()
}

// reverse returns input string reversed
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func listen(port string, registry *prometheus.Registry) {
	log.Println("Listening on port", port)
	router := chi.NewRouter()
	router.Get("/", ReturnRelease)
	router.Post("/", ReverseWord)
	router.Get("/hostname", ReturnHostname)
	router.Get("/health", ReturnHealth)
	router.Mount("/fullmetrics", promhttp.Handler()) // Default prometheus collector, includes other stuff on top of our custom registers
	router.Mount("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// getEnv returns the value for a given Env Var
func getEnv(varName, defaultValue string) string {
	if varValue, ok := os.LookupEnv(varName); ok {
		return varValue
	}
	return defaultValue
}

func main() {
	release := getEnv("RELEASE", "NotSet")
	port := getEnv("APP_PORT", "8080")
	log.Println("Starting Reverse Api", version, "Release:", release)
	// Custom registry, this will be used by the /metrics endpoint and will only show the app metrics
	registry := prometheus.NewRegistry()
	// Add our custom registers to our custom registry
	registry.MustRegister(totalWordsReversed, endpointsAccessed)
	// Add our custom registers to the default register, this will be used by the /fullmetrics endpoints
	prometheus.MustRegister(totalWordsReversed, endpointsAccessed)
	// Remove the go/process collectors from the default prometheus registry, this can be used to remove collectos from the default prometheus registry
	//prometheus.Unregister(collectors.NewGoCollector())
	//prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	listen(port, registry)
}
