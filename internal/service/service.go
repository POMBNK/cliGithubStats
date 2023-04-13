package service

import (
	"log"
	"os/user"

	"github.com/POMBNK/cliGitStats/pkg/scanHelpers/finder"
	"github.com/POMBNK/cliGitStats/pkg/scanHelpers/writer"
	"github.com/POMBNK/cliGitStats/pkg/statHelpers/collector"
	"github.com/POMBNK/cliGitStats/pkg/statHelpers/printer"
)

const dotfileDirectory = "/.gogitstats"

type Service struct {
	folder    string
	email     string
	finder    *finder.Finder
	writer    *writer.Writer
	collector *collector.Collector
	printer   *printer.Printer
}

// Service constructor
func New(folder string, email string, finder *finder.Finder, writer *writer.Writer, collector *collector.Collector, printer *printer.Printer) *Service {
	return &Service{
		folder:    folder,
		email:     email,
		finder:    finder,
		writer:    writer,
		collector: collector,
		printer:   printer,
	}
}

// Run start the application
func (s *Service) Run() {

	if err := s.findFiles(); err != nil {
		log.Fatalf("Can't run service findFiles() %v", err)
	}

	if err := s.showStat(); err != nil {
		log.Fatalf("Can't run service showStats() %v", err)
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

func (s *Service) getStats() (map[int]int, error) {
	commits, err := s.collector.ProcessRepos(s.email)
	if err != nil {
		return nil, err
	}

	return commits, nil
}

func (s *Service) showStat() error {
	commits, err := s.getStats()
	if err != nil {
		return err
	}
	s.printer.Show(commits)

	return nil
}
