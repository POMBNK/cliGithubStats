package finder

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	gitDirectory     = "/.git"
	gitDirectoryName = ".git"
)

type Finder struct {
}

func New() *Finder {
	return &Finder{}
}

// RecursiveSearch return FindGitDirs.
func (fd *Finder) RecursiveSearch(folder string) []string {
	return fd.findGitDirs(make([]string, 0), folder)
}

// FindGitDirs recursively find .git folders paths in given folder directory.
// All folder paths added to slice.
func (fd *Finder) findGitDirs(folders []string, baseDir string) []string {

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
			path = baseDir + "/" + file.Name()
			if file.Name() == gitDirectoryName {
				path = strings.TrimSuffix(path, gitDirectory)
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}
			folders = fd.findGitDirs(folders, path)
		}
	}

	return folders
}
