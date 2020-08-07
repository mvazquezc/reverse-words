package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
			Name: "total_reversed_words",
			Help: "Total number of reversed words",
		},
	)
)

var (
	endpointsAccessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "endpoints_accessed",
			Help: "Total number of accessed to a given endpoint",
		},
		[]string{"accessed_endpoint"},
	)
)

var version = "v0.0.15"

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

//ReturnHealth returns healthy string, can be used for monitoring pourposes
func ReturnHealth(w http.ResponseWriter, r *http.Request) {
	health := "Healthy"
	_, err := w.Write([]byte(health + "\n"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	endpointsAccessed.WithLabelValues("health").Inc()
}

//ReverseWord returns a reversed word based on an input word
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

//reverse returns input string reversed
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

//getEnv returns the value for a given Env Var
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
	log.Println("Listening on port", port)
	prometheus.MustRegister(totalWordsReversed)
	prometheus.MustRegister(endpointsAccessed)
	router := mux.NewRouter()
	router.HandleFunc("/", ReverseWord).Methods("POST")
	router.HandleFunc("/", ReturnRelease).Methods("GET")
	router.HandleFunc("/hostname", ReturnHostname).Methods("GET")
	router.HandleFunc("/health", ReturnHealth).Methods("GET")
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
