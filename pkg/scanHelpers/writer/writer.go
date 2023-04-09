package writer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Writer struct {
}

func New() *Writer {
	return &Writer{}
}

// ParseFile  parse dotfile and returning slice of parsed lines
func (w *Writer) ParseFile(filePath string) ([]string, error) {
	f, err := w.openFile(filePath)
	defer f.Close()

	if err != nil {
		return nil, fmt.Errorf("can't open file %s", err)
	}

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			return nil, err
		}
	}

	return lines, nil
}

// OpenFile If dotfile doesn't exist, create a new dotfile else open it with 0755 permission.
func (w *Writer) openFile(filepath string) (*os.File, error) {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(filepath)
			if err != nil {
				return nil, fmt.Errorf("can't create current file %w", err)
			}
		} else {
			return nil, fmt.Errorf("something went wrong with opennig file %w", err)
		}
	}
	return f, nil
}

// JoinRow adds new elements to existing slice only if this element not already there
func (w *Writer) JoinRow(new []string, existing []string) []string {
	for _, row := range new {
		if !contains(existing, row) {
			existing = append(existing, row)
		}
	}

	return existing
}

// WriteToFile writes found repos to dotfile.
func (w *Writer) WriteToFile(repos []string, filepath string) error {
	content := strings.Join(repos, "\n")
	if err := os.WriteFile(filepath, []byte(content), 0755); err != nil {
		return nil
	}
	return nil
}

// contains checks is value already in slice.
func contains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}

	return false
}
