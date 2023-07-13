package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
)

var port = flag.String("port", "32778", "TCP port to bind to")
var directory = flag.String("directory", "images", "directory containing images")
var path = flag.String("path", "photo", "path to serve image on")

func main() {
	flag.Parse()

	d, err := filepath.Abs(*directory)
	if err != nil {
		panic(err)
	}

	log.Printf("Serving images from %v at http://localhost:%s/%s", d, *port, *path)

	handler := http.HandlerFunc(handleRequest)
	http.Handle("/"+*path, handler)
	http.ListenAndServe(":"+*port, nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
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
