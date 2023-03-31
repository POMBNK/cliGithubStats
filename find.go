package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
)

const (
	gitDirectoryName = ".git"
	gitDirectory     = "/.git"
	dotfileDirectory = "/.gogitstats"
	home             = "path_to_folder"
)

func main() {
	err := Scan(home)
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func Scan(folder string) error {
	fmt.Println("Loading...")
	repos := RecursiveSearch(folder)
	dotfile := GetDotFilePath()
	if err := ChangeFile(dotfile, repos); err != nil {
		return err
	}
	fmt.Printf("\nSuccess!\n")

	return nil
}

// RecursiveSearch return FindGitDirs.
func RecursiveSearch(folder string) []string {
	return FindGitDirs(make([]string, 0), folder)
}

// FindGitDirs recursively find .git folders paths in given folder directory.
// All folder paths added to slice.
func FindGitDirs(folders []string, baseDir string) []string {

	//trim last "/" in path
	baseDir = strings.TrimSuffix(baseDir, "/")

	f, err := os.Open(baseDir)
	if err != nil {
		log.Fatalf("can't open file %s", err)
	}

	files, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		log.Fatalf("can't get files in some directory %s", err)
	}

	var path string

	// Check all found directories inside them
	for _, file := range files {
		if file.IsDir() {
			fileName := file.Name()
			path = baseDir + "/" + fileName
			if fileName == gitDirectoryName {
				path = strings.TrimSuffix(path, gitDirectory)
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}
			folders = FindGitDirs(folders, path)
		}
	}

	return folders
}

// GetDotFilePath Get path to dotfile.
func GetDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Can't get current User information %s", err)
	}
	dotFile := usr.HomeDir + dotfileDirectory

	return dotFile
}

/*-----------------------------------------------Work with dotfile----------------------------------------------------*/

func ChangeFile(filepath string, newRepos []string) error {
	existingRepos, err := ParseFile(filepath)
	if err != nil {
		return err
	}
	repos := JoinRow(newRepos, existingRepos)
	if err = WriteToFile(repos, filepath); err != nil {
		return err
	}
	return nil
}

// ParseFile  parse dotfile and returning slice of parsed lines
func ParseFile(filePath string) ([]string, error) {
	f, err := OpenFile(filePath)
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

// OpenFile If dotfile doesn't exist, create a new dotfile else open it with 0755 permission.
func OpenFile(filepath string) (*os.File, error) {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0755)
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

func JoinRow(new []string, existing []string) []string {
	for _, row := range new {
		if !contains(existing, row) {
			existing = append(existing, row)
		}
	}

	return existing
}

func contains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}

	return false
}

func WriteToFile(repos []string, filepath string) error {
	content := strings.Join(repos, "\n")
	if err := os.WriteFile(filepath, []byte(content), 0755); err != nil {
		return nil
	}
	return nil
}
