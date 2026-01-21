# 🎓 Learning & Interview Guide

This document helps you synthesize knowledge from the **ImageToExcel Importer** project and prepare for Golang/Backend interviews.

## 🧠 Core Concepts

### 1. Concurrency
The project uses the **Worker Pool** pattern to process images.
- **Why Worker Pool?**: To control the number of concurrent Goroutines. If you create one Goroutine for every image file (e.g., 10,000 images), the system will suffer from thrashing and run out of RAM.
- **Channels**: Uses a `jobs channel` to send tasks and a `results channel` to receive results. This is the "Fan-out / Fan-in" pattern.
- **Synchronization**: Uses `sync.WaitGroup` to wait for all workers to complete before closing the result channel.

### 2. Memory Management
- **Streaming Excel**: Uses `rows.Next()` from the `excelize` library instead of reading the entire file into RAM. This allows processing Excel files with millions of rows while keeping memory usage stable.
- **Lazy Loading**: Only loads the image decode config (`image.DecodeConfig`) to get dimensions before loading the full pixel data.

### 3. Application Architecture
- **Wails Framework**: Combines the power of Go (Backend performance) and Web Tech (Frontend UI).
- **Frontend-Backend Bridge**: Asynchronous communication via a JSON bridge.

---

## 🎤 Interview Questions

Here are questions recruiters might ask based on this project:

### Level: Junior / Fresher

**Q1: Why did you choose Go for this project instead of Python or C#?**
> *Hint:* Go has fast startup times, compiles to a compact native binary (no heavy runtime like .NET/JVM required), and especially its Concurrency model (Goroutines) is very powerful for handling I/O bound tasks (reading/writing image files) and CPU bound tasks (image compression) simultaneously.

**Q2: How do you ensure thread-safety when writing to the Excel file?**
> *Hint:* The `excelize` library is not strictly thread-safe for parallel writes. In this project, I utilize the "Single Consumer" pattern: Multiple workers process images in parallel, but results are pushed to a single channel. This channel is consumed by a single loop (main thread) that writes to Excel sequentially. This eliminates race conditions without needing complex Mutexes.

**Q3: How does `defer` work? Why use `defer wg.Done()`?**
> *Hint:* `defer` pushes a function onto a stack to be executed in LIFO order when the surrounding function returns. Using `defer wg.Done()` ensures that even if a worker panics or returns early, the `WaitGroup` counter is still decremented, preventing deadlocks (program hanging forever).

### Level: Mid / Senior

**Q4: How do you handle UI updates (Progress Bar) from the Go Backend in Wails?**
> *Hint:* Wails provides an `EventsEmit` mechanism. From the Go backend, I emit an event containing the percentage progress. The Frontend listens for this event (`runtime.EventsOn`) and updates the DOM. This is an Event-Driven model, helping decouple backend logic from the interface.

**Q5: If the Excel file has 1 million rows, can the current project handle it? Will it OOM (Out Of Memory)?**
> *Hint:* It can handle it thanks to using the `Iterator` (`rows.Next()`) to fetch one row at a time instead of `GetRows()` which loads the whole block. However, the `productMap` storing product codes still resides in RAM. With 1 million rows, this map might consume tens to hundreds of MB of RAM, which is still within acceptable limits for modern computers. For further optimization, a lightweight embedded database (SQLite/BadgerDB) could replace the in-memory map.

**Q6: How do you optimize Docker/CI build speed for a Go project?**
> *Hint:* Use Cache for `go mod download` and `go build` (as configured in GitHub Actions `setup-go` with `cache: true`). Use multi-stage builds in the Dockerfile (build in stage 1, copy binary to alpine/scratch in stage 2) to reduce image size.

**Q7: How do you design the Auto-Update feature securely?**
> *Hint:*
> 1. Check checksum/hash of the downloaded file (currently a point for improvement).
> 2. Use atomic replacement mechanisms or a batch script buffer to avoid "file in use" errors on Windows.
> 3. Clear Versioning (Semantic Versioning) and inject version at build time using `ldflags` to avoid hardcoding errors.

---

## 📚 Extended Challenges

To master the concepts, try implementing the following tasks yourself:
1. **Add Checksum Validation**: Verify the SHA256 hash of the downloaded update against the file on GitHub Release.
2. **Stop/Resume**: Add a "Pause" button to temporarily stop the worker pool and a "Resume" button to continue.
3. **Benchmarking**: Write benchmarks comparing processing speed when using `WorkerCount = 1` vs `WorkerCount = 10`.
