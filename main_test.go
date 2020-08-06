package main

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "strings"
    "os"
)

func TestReturnRelease(t *testing.T) {
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(ReturnRelease)
    handler.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK {
        t.Errorf("ReturnRelease returned wrong status code, expected %v, got %v", http.StatusOK, rr.Code)
    }
    expectedRelease := "NotSet"
    if ! strings.Contains(rr.Body.String(), expectedRelease) {
        t.Errorf("ReturnRelease returned wrong release, expected %v, got %v", expectedRelease, rr.Body.String())
    }
}

func TestReturnHostname(t *testing.T) {
    req, err := http.NewRequest("GET", "/hostname", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(ReturnHostname)
    handler.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK {
        t.Errorf("ReturnHostname returned wrong status code, expected %v, got %v", http.StatusOK, rr.Code)
    }
    expectedHostname, err := os.Hostname()
    if err != nil {
        t.Fatal(err)
    }

    if ! strings.Contains(rr.Body.String(), expectedHostname) {
        t.Errorf("ReturnHostname returned wrong hostname, expected %v, got %v", expectedHostname, rr.Body.String())
    }
}

func TestReturnHealth(t *testing.T) {
    req, err := http.NewRequest("GET", "/health", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(ReturnHealth)
    handler.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK {
        t.Errorf("ReturnHealth returned wrong status code, expected %v, got %v", http.StatusOK, rr.Code)
    }
    expectedHealth := "Healthy"

    if ! strings.Contains(rr.Body.String(), expectedHealth) {
        t.Errorf("ReturnHealth returned wrong health status, expected %v, got %v", expectedHealth, rr.Body.String())
    }
}

func TestReverseWord(t *testing.T) {
    body := strings.NewReader(`{"word": "PALC"}`)
    req, err := http.NewRequest("POST", "/", body)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(ReverseWord)
    handler.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK {
        t.Errorf("ReverseWord returned wrong status code, expected %v, got %v", http.StatusOK, rr.Code)
    }
    expectedResponse := `{"reverse_word":"CLAP"}`

    if ! strings.Contains(rr.Body.String(), expectedResponse) {
        t.Errorf("ReverseWord returned wrong word, expected %v, got %v", expectedResponse, rr.Body.String())
    }
}

func TestGetEnv(t *testing.T) {
    os.Setenv("FOO", "BAR")
    result := getEnv("FOO", "DEFAULT_VALUE")
    if result != "BAR" {
        t.Errorf("TestGetEnv with existing variable failed, expected %v, got %v", "BAR", result)
    }
    result = getEnv("NOT_EXISTING_VAR", "DEFAULT_VALUE")
    if result != "DEFAULT_VALUE" {
        t.Errorf("TestGetEnv with default value failed, expected %v, got %v", "DEFAULT_VALUE", result)
    }
}

func TestReverse(t *testing.T) {
    result := reverse("PALC")
    if result != "CLAP" {
        t.Errorf("TestReverse failed, expected %v, got %v", "CLAP", result)
    }
}