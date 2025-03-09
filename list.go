package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/djherbis/times"
)

type Cameras struct {
	Cameras []string `json:"cameras"`
}

func list_directories_handler(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir(lib_directory)
	if err != nil {
		fmt.Fprintf(w, "can't read lib directory.")
	}

	var directories []string
	for idx := 0; idx < len(entries); idx++ {
		if entries[idx].IsDir() {
			directories = append(directories, entries[idx].Name())
		}
	}
	res := Cameras{Cameras: directories}
	json.NewEncoder(w).Encode(res)
}

type Images struct {
	Camera string  `json:"camera"`
	Images []Image `json:"images"`
}

type Image struct {
	Name       string `json:"name"`
	Created_at string `json:"created_at"`
}

func list_images_handler(w http.ResponseWriter, r *http.Request) {
	camera := r.URL.Query().Get("camera")
	image_directory := lib_directory + "/" + camera

	entries, err := os.ReadDir(image_directory)
	if err != nil {
		fmt.Printf("can't read lib directory at %s", image_directory)
		fmt.Fprintf(w, "can't read lib directory.")
	}

	var images []Image
	for idx := 0; idx < len(entries); idx++ {
		file := entries[idx]
		filename := file.Name()
		if strings.HasPrefix(filename, ".") {
			continue
		}
		if !file.Type().IsRegular() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			fmt.Println("warning: file has no info()")
		}

		ts, err := times.Stat(image_directory + "/" + file.Name())
		if err != nil {
			fmt.Printf("can't call .Stat on %s: %d\n", file.Name(), err)
		}

		file_creation_ts := info.ModTime().String()
		if ts.HasBirthTime() {
			file_creation_ts = ts.BirthTime().String()
		}

		images = append(images, Image{Name: filename, Created_at: file_creation_ts})
	}
	res := Images{Camera: camera, Images: images}
	json.NewEncoder(w).Encode(res)
}
