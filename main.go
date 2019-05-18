package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
        "os"
)

type Input struct {
	Word string `json:"word,omitempty"`
}
type Output struct {
	ReverseWord string `json:"reverse_word,omitempty"`
}

func ReturnRelease(w http.ResponseWriter, r *http.Request) {
        release := getEnv("RELEASE", "NotSet")
        releaseString := "Reverse Words Release: " + release
        w.Write([]byte(releaseString))
}

func ReturnHealth(w http.ResponseWriter, r *http.Request) {
        health := "Healthy"
        w.Write([]byte(health))
}

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
	output := Output{reverseWord}
        js, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func reverse(s string) string {
  runes := []rune(s)
  for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
    runes[i], runes[j] = runes[j], runes[i]
  }
  return string(runes)
}

func getEnv(varName, defaultValue string) string {
    if varValue, ok := os.LookupEnv(varName); ok {
        return varValue
    }
    return defaultValue
}

func main() {
        release := getEnv("RELEASE", "NotSet")
        port := getEnv("APP_PORT", "8080")
	log.Println("Starting Reverse Api. Release:", release)
	log.Println("Listening on port", port)
	router := mux.NewRouter()
	router.HandleFunc("/", ReverseWord).Methods("POST")
        router.HandleFunc("/", ReturnRelease).Methods("GET")
        router.HandleFunc("/health", ReturnHealth).Methods("GET")
	log.Fatal(http.ListenAndServe(":" + port, router))
}
