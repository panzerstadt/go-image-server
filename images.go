package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
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

type Size struct {
	name   string
	height int
	width  int
}

func get_or_compute_size(filepath string) Size {
	cached := from_cache(filepath)
	if cached.name != "" {
		return cached
	}

	image_file, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("can't load image", http.StatusGone)
	}

	image_file_decoded, err := jpeg.Decode(bytes.NewReader(image_file))
	if err != nil {
		fmt.Printf("image could not be decoded: %s", filepath)
	}
	image_size := image_file_decoded.Bounds()
	width := image_size.Dx()
	height := image_size.Dy()

	result := Size{name: filepath, width: width, height: height}
	to_cache(result)

	return result
}
