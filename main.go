package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	lib_directory string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	lib_directory = os.Getenv("LIB_DIRECTORY")
}

func logger(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		route := r.URL.Path
		params := r.URL.Query()
		handler(w, r)
		elapsed := time.Since(start)
		fmt.Printf("%s %12s %s %s\n", route, elapsed, time.Now().Local().Format("2006-01-02T15:04:05"), params)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "go away")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	http.HandleFunc("/favicon.ico", logger(handler))
	http.HandleFunc("/", logger(list_directories_handler))
	http.HandleFunc("/list", logger(list_directories_handler))
	http.HandleFunc("/camera", logger(list_images_handler))
	http.HandleFunc("/images", logger(serve_images))
	http.HandleFunc("/healthcheck", logger(healthcheck))

	fmt.Println("serving images at :8080")
	monitor(ctx)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("server encountered an error and will shut down")
		cancel()
	}
	fmt.Println("shutting down server.")
}
