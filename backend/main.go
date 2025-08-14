package main

import (
	"fmt"
	"os"

	"github.com/LaurelEdison/clashcoder/backend/server"
)

func main() {
	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting server: %v", err)
		os.Exit(1)
	}
}
