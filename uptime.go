package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

func monitor() {
	for {
		go_check()
		proc_check()
		time.Sleep(10 * time.Second)
	}
}

func go_check() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	fmt.Printf("cpu: %d, goroutines: %d, mem: %d bytes, mallocs: %d\n", runtime.NumCPU(), runtime.NumGoroutine(), mem.TotalAlloc, mem.Mallocs)
}

func proc_check() {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err == nil {
		fmt.Println(stat)
	}
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
