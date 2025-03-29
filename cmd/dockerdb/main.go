package main

import (
	"dockerdb/internal/cli"
	// Use an alias for one of the cli packages
)

func main() {
	cli.Execute()
	// Parse command-line arguments
	// if err := app.Run(os.Args); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	// 	os.Exit(1)
	// }
}
