package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var port = flag.String("port", "32778", "TCP port to bind to")
var directory = flag.String("directory", "images", "directory containing images")

func main() {
	flag.Parse()

	d, err := filepath.Abs(*directory)
	if err != nil {
		panic(err)
	}

	log.Printf("Serving images from %v at http://localhost:%s/photo", d, *port)

	mux := http.NewServeMux()
	mux.HandleFunc("/photo", getPhoto)
	mux.HandleFunc("/outdated", getOutdated)
	http.ListenAndServe(":"+*port, mux)
}

func getPhoto(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(*directory)
	if err != nil {
		panic(err)
	}
	files, err := os.ReadDir(*directory)
	if err != nil {
		panic(err)
	}
	file := files[rand.Intn(len(files))]
	fileBytes, err := os.ReadFile(filepath.Join(path, file.Name()))
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}

func getOutdated(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(*directory)
	if err != nil {
		panic(err)
	}
	files, err := os.ReadDir(*directory)
	if err != nil {
		panic(err)
	}
	file := files[rand.Intn(len(files))]
	fileBytes, err := os.ReadFile(filepath.Join(path, file.Name()))
	if err != nil {
		panic(err)
	}
	// Create a last modified date in the past, so that "helpful" browsers will
	// serve the image from the cache.
	thePast := time.Now().Add(-15 * time.Minute).UTC().Format(http.TimeFormat)
	w.Header().Set("Last-Modified", thePast)
	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(fileBytes)
}
