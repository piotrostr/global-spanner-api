package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckWorks(t *testing.T) {
	client, err := SetupClient()
	if err != nil {
		t.Errorf("Error setting up client: %v", err)
	}
	router := SetupRouter(client)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf(
			"Expected status code %d, got %d: %s",
			http.StatusCreated,
			w.Code,
			w.Body.String(),
		)
	}
}

func TestAddNamesWorks(t *testing.T) {
	client, err := SetupClient()
	if err != nil {
		t.Errorf("Error setting up client: %v", err)
	}
	router := SetupRouter(client)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/add-names", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf(
			"Expected status code %d, got %d: %s",
			http.StatusCreated,
			w.Code,
			w.Body.String(),
		)
	}
}

func TestGetNamesWorks(t *testing.T) {
	client, err := SetupClient()
	if err != nil {
		t.Errorf("Error setting up client: %v", err)
	}
	router := SetupRouter(client)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/get-names", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf(
			"Expected status code %d, got %d: %s",
			http.StatusCreated,
			w.Code,
			w.Body.String(),
		)
	}
}
