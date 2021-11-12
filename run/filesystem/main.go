package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const mntDir = "/mnt/nfs/filestore"
const filename = "test-Wed-Feb-23-16:48-2022.txt"

func handler(w http.ResponseWriter, r *http.Request) {
  path := filepath.Join(mntDir, filename)
  // fi, err := os.Stat(path)
  // if err != nil {
  //   io.WriteString(w, fmt.Sprintf("%s: %v\n", path, err))
  //   return
  // }

  contents, err := ioutil.ReadFile(path)
  if err != nil {
    io.WriteString(w, fmt.Sprintf("Error retrieving file, %s: %v\n", path, err))
    return
  }
  // io.WriteString(w, string(contents))
	fmt.Fprint(w, string(contents))
}

func main() {
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: http.HandlerFunc(handler),
	}

	// Create channel to listen for signals.
	signalChan := make(chan os.Signal, 1)
	// SIGINT handles Ctrl+C locally.
	// SIGTERM handles Cloud Run termination signal.
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server.
	go func() {
		log.Printf("listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Receive output from signalChan.
	sig := <-signalChan
	log.Printf("%s signal caught", sig)

	// Timeout if waiting for connections to return idle.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add extra handling here to clean up resources, such as flushing logs and
	// closing any database or Redis connections.

	// Gracefully shutdown the server by waiting on existing requests (except websockets).
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %+v", err)
	}
	log.Print("server exited")
}