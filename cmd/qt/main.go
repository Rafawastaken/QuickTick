package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rafawastaken/quicktick/internal/cli"
)

func main() {
	_ = godotenv.Load() // Ignore error if .env file not found

	if err := cli.NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
