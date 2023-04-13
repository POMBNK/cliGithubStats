package collector

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	outOfRange          = 99999
	daysInLastSixMonths = 183
	dotfileDirectory    = "/.gogitstats"
)

type Collector struct {
}

func New() *Collector {
	return &Collector{}
}

func (c *Collector) ProcessRepos(email string) (map[int]int, error) {
	filePath := c.getDotFilePath()
	repos, err := c.parseFile(filePath)
	if err != nil {
		return nil, err
	}

	daysInMap := daysInLastSixMonths
	commits := make(map[int]int, daysInMap)
	for i := daysInMap; i > 0; i-- {
		commits[i] = 0
	}

	for _, path := range repos {
		commits = c.fillCommits(email, path, commits)
	}

	return commits, nil
}

func (c *Collector) fillCommits(email string, path string, commits map[int]int) map[int]int {
	// instantiate a git repo object from path
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}
	// get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}
	// get the commits history starting from HEAD
	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}
	// iterate the commits
	offset := c.calcOffset()
	err = iterator.ForEach(func(com *object.Commit) error {
		daysAgo := c.countDaysSinceDate(com.Author.When) + offset

		if com.Author.Email != email {
			return nil
		}

		if daysAgo != outOfRange {
			commits[daysAgo]++
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
	return commits
}

func (c *Collector) getBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}

func (c *Collector) countDaysSinceDate(date time.Time) int {
	days := 0
	now := c.getBeginningOfDay(time.Now())
	for date.Before(now) {
		date = date.Add(time.Hour * 24)
		days++
		if days > daysInLastSixMonths {
			return outOfRange
		}
	}
	return days
}

func (c *Collector) calcOffset() int {
	var offset int
	weekday := time.Now().Weekday()

	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1
	}

	return offset
}
