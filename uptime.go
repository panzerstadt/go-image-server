package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

// Add storage for previous measurements
var (
	prevCPUStats    *linuxproc.Stat
	prevCPUStatTime time.Time
)

const CHECK_FREQUENCY_SECONDS = 60

// check either every minute or everytime we serve images
func monitor() {
	go func() {
		for {
			go_check()
			proc_check()
			time.Sleep(CHECK_FREQUENCY_SECONDS * time.Second)
		}
	}()
}

func go_check() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	mb := float64(mem.TotalAlloc) / 1024 / 1024

	fmt.Printf("cpu: %d, goroutines: %d, mem: %.2f mb, mallocs: %d\n", runtime.NumCPU(), runtime.NumGoroutine(), mb, mem.Mallocs)
}

func proc_check() {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		return
	}
	meminfo, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		return
	}

	now := time.Now()

	// If we have previous measurements, calculate the difference
	if prevCPUStats != nil {
		// Calculate total CPU time difference
		prevTotal := prevCPUStats.CPUStatAll.User + prevCPUStats.CPUStatAll.Nice +
			prevCPUStats.CPUStatAll.System + prevCPUStats.CPUStatAll.Idle +
			prevCPUStats.CPUStatAll.IOWait + prevCPUStats.CPUStatAll.IRQ +
			prevCPUStats.CPUStatAll.SoftIRQ + prevCPUStats.CPUStatAll.Steal +
			prevCPUStats.CPUStatAll.Guest

		currentTotal := stat.CPUStatAll.User + stat.CPUStatAll.Nice +
			stat.CPUStatAll.System + stat.CPUStatAll.Idle +
			stat.CPUStatAll.IOWait + stat.CPUStatAll.IRQ +
			stat.CPUStatAll.SoftIRQ + stat.CPUStatAll.Steal +
			stat.CPUStatAll.Guest

		totalDiff := currentTotal - prevTotal

		// Calculate idle time difference
		prevIdle := prevCPUStats.CPUStatAll.Idle + prevCPUStats.CPUStatAll.IOWait
		currentIdle := stat.CPUStatAll.Idle + stat.CPUStatAll.IOWait
		idleDiff := currentIdle - prevIdle

		// Calculate CPU usage percentage
		cpuUsage := 100 * (1.0 - float64(idleDiff)/float64(totalDiff))

		// Calculate per-core CPU usage
		var coreUsages []float64
		for i, s := range stat.CPUStats {
			if i >= len(prevCPUStats.CPUStats) {
				break
			}

			prevCore := prevCPUStats.CPUStats[i]
			prevTotal := prevCore.User + prevCore.Nice + prevCore.System +
				prevCore.Idle + prevCore.IOWait + prevCore.IRQ +
				prevCore.SoftIRQ + prevCore.Steal + prevCore.Guest

			currentTotal := s.User + s.Nice + s.System + s.Idle +
				s.IOWait + s.IRQ + s.SoftIRQ + s.Steal + s.Guest

			totalDiff := currentTotal - prevTotal

			prevIdle := prevCore.Idle + prevCore.IOWait
			currentIdle := s.Idle + s.IOWait
			idleDiff := currentIdle - prevIdle

			if totalDiff > 0 {
				coreUsage := 100 * (1.0 - float64(idleDiff)/float64(totalDiff))
				coreUsages = append(coreUsages, coreUsage)
			} else {
				coreUsages = append(coreUsages, 0.0)
			}
		}

		// Calculate memory usage percentage
		memUsed := meminfo.MemTotal - meminfo.MemFree - meminfo.Buffers - meminfo.Cached
		memUsagePercent := (float64(memUsed) / float64(meminfo.MemTotal)) * 100

		if len(coreUsages) >= 4 {
			fmt.Printf("cpu: %.2f%% [cores: %.2f%%, %.2f%%, %.2f%%, %.2f%%]\n",
				cpuUsage, coreUsages[0], coreUsages[1], coreUsages[2], coreUsages[3])
		} else {
			fmt.Printf("cpu: %.2f%% [cores: %v]\n", cpuUsage, coreUsages)
		}
		fmt.Printf("mem: %.2f%%\n", memUsagePercent)
	}

	// Store current values for next iteration
	prevCPUStats = stat
	prevCPUStatTime = now
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
