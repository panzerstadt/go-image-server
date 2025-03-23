package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var cacheMutex sync.RWMutex

func check_cache_size() {
	f, _ := get_cache()
	defer f.Close()
	stat, _ := f.Stat()
	size := stat.Size()
	fmt.Printf("cache size: %d mb", size/1024/1024)
}

func get_cache() (*os.File, error) {
	file, err := os.OpenFile("sizes", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		file.Close()
		return nil, err
	}
	return file, err
}

func to_cache(size Size) error {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	file, err := get_cache()
	if err != nil {
		return err
	}
	defer file.Close()

	file.SetWriteDeadline(time.Now().Add(1 * time.Second))
	written, err := file.Write([]byte(fmt.Sprintf("\n%s %d %d", size.name, size.height, size.width)))
	if err != nil {
		return err
	}
	if written == 0 {
		return fmt.Errorf("caching failed: %v", size)
	}
	check_cache_size()
	return nil
}

func from_cache(key string) Size {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	file, err := get_cache()
	if err != nil {
		fmt.Println("error loading cache")
		return Size{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, key) {
			tokens := strings.Split(line, " ")
			height, err := strconv.Atoi(tokens[1])
			if err != nil {
				fmt.Println("Error parsing height:", err)
				return Size{}
			}
			width, err := strconv.Atoi(tokens[2])
			if err != nil {
				fmt.Println("Error parsing width:", err)
				return Size{}
			}
			return Size{name: tokens[0], height: height, width: width}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading sizes file:", err)
		return Size{}
	}

	return Size{}
}
