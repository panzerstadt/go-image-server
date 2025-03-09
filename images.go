package main

import (
	"net/http"
	"os"
	"strconv"
)

func serve_images(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		http.Error(w, "image name not found", http.StatusNotFound)
		return
	}

	image_file, err := os.ReadFile(lib_directory + "/" + name)
	if err != nil {
		http.Error(w, "can't load image", http.StatusGone)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(image_file)))
	_, err = w.Write(image_file)
	if err != nil {
		http.Error(w, "error sending file", http.StatusInternalServerError)
	}
}
