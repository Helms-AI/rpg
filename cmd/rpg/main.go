// Package main is the entry point for the rpg MCP server.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kon1790/rpg/internal/server"
)

func main() {
	// Parse command line flags
	outputDir := flag.String("output", "./output", "Base output directory for generated projects (language subdirs will be created)")
	flag.StringVar(outputDir, "o", "./output", "Base output directory (shorthand)")
	flag.Parse()

	// Create context that cancels on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Create and run the MCP server
	srv := server.New(*outputDir)
	if err := srv.Run(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
