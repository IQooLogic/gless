package reader

import (
	"bufio"
	"errors"
	"io"
	"os"
)

// FileReader reads a file line by line with buffering
type FileReader struct {
	file     *io.ReadCloser
	lines    []string
	loaded   bool
	filename string
}

// NewFileReader creates a new file reader
func NewFileReader(filename string) (*FileReader, error) {
	var file io.ReadCloser

	if filename == "-" || filename == "" {
		file = os.Stdin
		filename = "stdin"
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		file = f
	}

	return &FileReader{
		file:     &file,
		lines:    make([]string, 0),
		loaded:   false,
		filename: filename,
	}, nil
}

// Load reads all lines from the file into memory
func (fr *FileReader) Load() error {
	if fr.loaded {
		return nil
	}

	scanner := bufio.NewScanner(*fr.file)
	// Increase buffer size for long lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024) // Max 1MB per line

	for scanner.Scan() {
		fr.lines = append(fr.lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fr.loaded = true
	return nil
}

// GetLine returns the line at the specified index (0-based)
func (fr *FileReader) GetLine(index int) (string, error) {
	if !fr.loaded {
		if err := fr.Load(); err != nil {
			return "", err
		}
	}

	if index < 0 || index >= len(fr.lines) {
		return "", errors.New("line index out of bounds")
	}

	return fr.lines[index], nil
}

// GetLines returns a range of lines [start, end)
func (fr *FileReader) GetLines(start, end int) ([]string, error) {
	if !fr.loaded {
		if err := fr.Load(); err != nil {
			return nil, err
		}
	}

	if start < 0 {
		start = 0
	}
	if end > len(fr.lines) {
		end = len(fr.lines)
	}
	if start >= end {
		return []string{}, nil
	}

	return fr.lines[start:end], nil
}

// LineCount returns the total number of lines
func (fr *FileReader) LineCount() int {
	if !fr.loaded {
		fr.Load() // Ignore error, will return 0
	}
	return len(fr.lines)
}

// Filename returns the name of the file being read
func (fr *FileReader) Filename() string {
	return fr.filename
}

// Close closes the underlying file
func (fr *FileReader) Close() error {
	if fr.file != nil {
		return (*fr.file).Close()
	}
	return nil
}
