package collector

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
)

func (c *Collector) getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dotFile := usr.HomeDir + dotfileDirectory

	return dotFile
}

func (c *Collector) parseFile(filePath string) ([]string, error) {
	f, err := c.openFile(filePath)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("can't open file %s", err)
	}
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err(); err != nil {
		if err != io.EOF {
			return nil, err
		}
	}
	return lines, nil
}

func (c *Collector) openFile(filepath string) (*os.File, error) {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(filepath)
			if err != nil {
				return nil, fmt.Errorf("can't create current file %s", err)
			}
		} else {
			return nil, fmt.Errorf("something went wrong with opennig file %s", err)
		}
	}

	return f, nil
}
