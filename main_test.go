package main

import (
	"encoding/json"
	"flag"
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
	prod := flag.Bool("prod", false, "Run tests against production database")
	createTable := flag.Bool("create-table", false, "Create table in production database")
	flag.Parse()

	if *prod {
		log.Println("Running tests against production database")
	}

	if !*prod {
		srv, err := spannertest.NewServer("localhost:0")
		if err != nil {
			log.Fatal(err)
		}
		defer srv.Close()

		// if the environment variable SPANNER_EMULATOR_HOST is set, the
		// client will connect to the emulator instead of the real Cloud
		// Spanner service.
		os.Setenv("SPANNER_EMULATOR_HOST", srv.Addr)
	}

	var err error
	client, err = SetupClient()
	if err != nil {
		log.Fatalf("Error setting up client: %v", err)
	}
	router = SetupRouter(client)

	log.Println("Spanner URL:", client.spannerURL)

	// create table if flagged or not if in test environment
	if *createTable || !*prod {
		err = client.CreateTable()
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		}
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

	var names []Name
	if err := json.Unmarshal(w.Body.Bytes(), &names); err != nil {
		t.Errorf("Error unmarshaling names: %v", err)
	}
}
