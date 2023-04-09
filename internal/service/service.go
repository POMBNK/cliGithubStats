package service

import (
	"log"
	"os/user"

	"github.com/POMBNK/cliGitStats/pkg/scanHelpers/finder"
	"github.com/POMBNK/cliGitStats/pkg/scanHelpers/writer"
)

const dotfileDirectory = "/.gogitstats"

type Service struct {
	folder string
	finder *finder.Finder
	writer *writer.Writer
}

// Service constructor
func New(folder string, finder *finder.Finder, writer *writer.Writer) *Service {
	return &Service{
		folder: folder,
		finder: finder,
		writer: writer,
	}
}

// Run start the application
func (s *Service) Run() {

	if err := s.findFiles(); err != nil {
		log.Fatalf("Can't run service findFiles() %v", err)
	}

}

// findFiles find all .git repos in base path and write all pathes to /.gogitstat file
func (s *Service) findFiles() error {
	repos := s.finder.RecursiveSearch(s.folder)
	dotfile := s.getDotFilePath()
	if err := s.changeFile(dotfile, repos); err != nil {
		return err
	}

	return nil
}

// changeFile write all pathes to /.gogitstat file
func (s *Service) changeFile(filepath string, newRepos []string) error {
	existingRepos, err := s.writer.ParseFile(filepath)
	if err != nil {
		return err
	}
	repos := s.writer.JoinRow(newRepos, existingRepos)
	if err := s.writer.WriteToFile(repos, filepath); err != nil {
		return err
	}
	return nil
}

// GetDotFilePath Get path to dotfile.
func (s *Service) getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Can't get current User information %s", err)
	}
	dotFile := usr.HomeDir + dotfileDirectory

	return dotFile
}
