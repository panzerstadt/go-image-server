package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

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
	Width      int    `json:"width"`
	Height     int    `json:"height"`
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
	ch := make(chan Image)
	var wg sync.WaitGroup

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

		wg.Add(1)
		go func(filename, image_directory, modified_at string) {
			defer wg.Done()
			get_image_stats(ch, filename, image_directory, modified_at)
		}(filename, image_directory, info.ModTime().String())
	}

	// Close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect results
	for data := range ch {
		images = append(images, data)
	}

	res := Images{Camera: camera, Images: images}
	json.NewEncoder(w).Encode(res)
}

func get_image_stats(ch chan<- Image, filename string, image_directory string, modified_at string) {
	image_filepath := image_directory + "/" + filename
	size := get_or_compute_size(image_filepath)

	ts, err := times.Stat(image_filepath)
	if err != nil {
		fmt.Printf("can't call .Stat on %s: %d\n", filename, err)
	}

	file_creation_ts := modified_at
	if ts.HasBirthTime() {
		file_creation_ts = ts.BirthTime().String()
	}

	ch <- Image{
		Name:       filename,
		Created_at: file_creation_ts,
		Width:      size.width,
		Height:     size.height,
	}
}
