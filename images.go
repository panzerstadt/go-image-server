package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
)

func serve_images(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		http.Error(w, "image name not found", http.StatusNotFound)
		return
	}

	image_file, err := os.Open(lib_directory + "/" + name)
	if err != nil {
		http.Error(w, "can't load image", http.StatusGone)
	}
	defer image_file.Close()

	w.Header().Set("Content-Type", "image/jpeg")
	_, err = io.Copy(w, image_file) // stream image directly
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

	img_config, err := jpeg.DecodeConfig(bytes.NewReader(image_file))
	if err != nil {
		fmt.Printf("image could not be decoded: %s", filepath)
	}
	result := Size{name: filepath, width: img_config.Width, height: img_config.Height}
	to_cache(result)

	return result
}
