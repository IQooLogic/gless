package main

import (
	"fmt"
	"os"

	"github.com/iqool/gless/internal/reader"
	"github.com/iqool/gless/internal/viewer"
)

func main() {
	// Parse command line arguments
	args := os.Args[1:]

	var filename string
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: gless <filename>")
		fmt.Fprintln(os.Stderr, "       gless -    (read from stdin)")
		os.Exit(1)
	}

	filename = args[0]

	// Create file reader
	fileReader, err := reader.NewFileReader(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer fileReader.Close()

	// Create and run viewer
	v := viewer.NewViewer(fileReader)
	if err := v.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running viewer: %v\n", err)
		os.Exit(1)
	}
}
