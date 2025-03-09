package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func get_cache() (*os.File, error) {
	file, err := os.OpenFile("sizes", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	return file, err
}

func to_cache(size Size) error {
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
	return nil
}

func from_cache(key string) Size {
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
