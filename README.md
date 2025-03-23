
### current issues
1. Unbounded goroutines:
    - Done: monitor() in uptime.go launches an infinite goroutine with no cancellation mechanism
    - Done: monitor_once() launches nested goroutines without proper tracking
2. File handling issues (ignored):
    - Done: get_cache() opens files but error cases might not properly close handles
    - The "sizes" cache file grows indefinitely with no cleanup mechanism
3. Image processing concerns:
    - Done: serve_images() loads entire images into memory at once -> switched to streaming
    - Done: get_or_compute_size() keeps full images in memory during processing
4. Resource management:
    - No timeout protection in list_images_handler() goroutines
    - Ignore: Cache has no eviction policy or size limits
    - Done: No mutual exclusion for cache file access (potential race conditions)

### learning notes for golang
- os.ReadFile() already calls defer file.close in its implementation, so i don't need to close it myself