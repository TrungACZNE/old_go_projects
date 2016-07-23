package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/TrungACZNE/git2go"
	"github.com/codegangsta/cli"
)

type Commit struct {
	Repo, Summary, Message string
	Date                   time.Time
}

type By func(p1, p2 *Commit) bool

func (by By) Sort(commits []Commit) {
	ps := &commitSorter{
		commits: commits,
		by:      by,
	}
	sort.Sort(ps)
}

type commitSorter struct {
	commits []Commit
	by      func(p1, p2 *Commit) bool
}

func (s *commitSorter) Len() int {
	return len(s.commits)
}

func (s *commitSorter) Swap(i, j int) {
	s.commits[i], s.commits[j] = s.commits[j], s.commits[i]
}

func (s *commitSorter) Less(i, j int) bool {
	return s.by(&s.commits[i], &s.commits[j])
}

func defaultSort(p1, p2 *Commit) bool {
	if p1.Repo < p2.Repo {
		return true
	}

	if p1.Repo > p2.Repo {
		return false
	}
	return p1.Date.After(p2.Date)
}

func start(repos, format string, weekDelta int) {
	commitHistory := make(map[int]map[string][]Commit)
	for _, repoPath := range strings.Split(repos, ",") {

		repo, err := git.OpenRepository(repoPath)
		if err != nil {
			log.Fatal(err)
		}
		defer repo.Free()

		walk, err := repo.Walk()
		if err != nil {
			log.Fatal(err)
		}
		defer walk.Free()

		err = walk.PushHead()
		if err != nil {
			log.Fatal(err)
		}

		walk.Iterate(func(commit *git.Commit) bool {
			commit, err := repo.LookupCommit(commit.Id())
			if err != nil {
				log.Fatal(err)
			}
			defer commit.Free()
			date := commit.Time()
			_, week := date.ISOWeek()
			email := commit.Author().Email

			var weekData map[string][]Commit
			var ok bool
			if weekData, ok = commitHistory[week]; !ok {
				weekData = make(map[string][]Commit)
				commitHistory[week] = weekData
			}

			if _, ok = weekData[email]; !ok {
				weekData[email] = []Commit{}
			}

			weekData[email] = append(weekData[email], Commit{
				Summary: commit.Summary(),
				Message: commit.Message(),
				Date:    date,
				Repo:    path.Base(repoPath),
			})

			return true
		})
	}

	_, thisWeek := time.Now().ISOWeek()
	commitSelectedWeek, ok := commitHistory[thisWeek-weekDelta]
	if !ok {
		log.Fatalln("No commit for selected week")
	}

	for email, commits := range commitSelectedWeek {
		if format == "brief" {
			fmt.Printf("Email \"%s\" has made %d commits\n", email, len(commits))
		} else {
			fmt.Printf("Email \"%s\" has made %d commits:\n", email, len(commits))
			By(defaultSort).Sort(commits)
			for _, commit := range commits {
				if format == "default" {
					fmt.Printf("\t%-30s %-30s %s\n", commit.Date.Format(time.RFC822), commit.Repo, commit.Summary)
				} else {
					fmt.Printf("\t%-30s %-30s\n\tSUMMARY: %s\n\tMESSAGE:\n%s\n", commit.Date.Format(time.RFC822), commit.Repo, commit.Summary, formatCommitMessage(commit.Message))
				}
			}
		}
	}
}

func formatCommitMessage(message string) string {
	message = strings.TrimSpace(message)
	result := ""
	for _, line := range strings.Split(message, "\n") {
		result += fmt.Sprintf("\t%s\n", line)
	}
	return result
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetFlags(0)

	app := cli.NewApp()
	app.Name = "Git weekly report"
	app.Usage = "Analyzes weekly commit history"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "repo",
			Value: "",
			Usage: "Comma separated list of repo directories",
		},
		cli.StringFlag{
			Name:  "format",
			Value: "default",
			Usage: "Sets report type - \"brief\": prints number of commits only, \"default\": prints number of commits and commit summaries, or \"verbose\": also prints commit messages",
		},
		cli.IntFlag{
			Name:  "weekDelta",
			Value: 0,
			Usage: "0 = this week, 1 = last week, etc",
		},
	}
	app.Action = func(c *cli.Context) {
		format := c.String("format")
		if format != "brief" && format != "default" && format != "verbose" {
			log.Fatal("Incorrect report format, please see --help")
		}
		start(c.String("repo"), format, c.Int("weekDelta"))
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Println("app.Run() error:", err)
	}
}
