package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cloud.google.com/go/spanner/spannertest"
	"github.com/gin-gonic/gin"
)

var client *Client

var router *gin.Engine

func TestMain(m *testing.M) {
	srv, err := spannertest.NewServer("localhost:0")
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Close()

	os.Setenv("SPANNER_EMULATOR_HOST", srv.Addr)

	client, err = SetupClient()
	if err != nil {
		log.Fatalf("Error setting up client: %v", err)
	}
	router = SetupRouter(client)

	err = client.CreateTable()
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	code := m.Run()

	client.Teardown()
	os.Exit(code)
}

func TestHealthCheckWorks(t *testing.T) {
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
	log.Println(w.Body.String())
}
