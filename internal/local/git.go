package local

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/DMcP89/tinycare-tui/internal/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetDailyCommits(path string) (string, error) {
	if path == "" {
		return "TINYCARE_WORKSPACE environment variable not set!", nil
	}
	repositories, err := FindGitRepositories(path)
	if err != nil {
		return "", fmt.Errorf("GetDailyCommits: Unable to find git repos for %s: %w", path, err)
	}

	result := ""
	for _, repo := range repositories {
		commitMessages, err := GetCommitsFromTimeRange(repo, time.Now().AddDate(0, 0, -1), time.Now())
		if err != nil {
			return "", err
		}

		if len(commitMessages) > 0 {
			result += fmt.Sprintf("[red]%s[white]\n", repo)
			result += commitMessages + "\n"
		}
	}

	return result, nil
}

func GetWeeklyCommits(path string) (string, error) {
	if path == "" {
		return "TINYCARE_WORKSPACE environment variable not set!", nil
	}
	repositories, err := FindGitRepositories(path)
	if err != nil {
		return "", err
	}

	result := ""
	for _, repo := range repositories {
		commitMessages, err := GetCommitsFromTimeRange(repo, time.Now().AddDate(0, 0, -7), time.Now())
		if err != nil {
			return "", err
		}

		if len(commitMessages) > 0 {
			result += fmt.Sprintf("[red]%s[white]\n", repo)
			result += commitMessages + "\n"
		}
	}

	return result, nil
}

func GetRepos(paths []string, c chan string, e chan error, q chan int) {
	var wg sync.WaitGroup
	wg.Add(len(paths))
	for _, path := range paths {
		go func(path string) {
			err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() && info.Name() == ".git" {
					c <- filepath.Dir(p)
				}
				return nil
			})
			if err != nil {
				e <- err
			}
			wg.Done()
		}(path)
	}
	wg.Wait()
	q <- 0
}

func FindGitRepositories(path string) ([]string, error) {
	var repositories []string
	//split the path into a slice of strings by comma
	repo_channel := make(chan string)
	error_channel := make(chan error)
	quit_channel := make(chan int)
	paths := strings.Split(path, ",")
	go GetRepos(paths, repo_channel, error_channel, quit_channel)
	for {
		select {
		case repo := <-repo_channel:
			repositories = append(repositories, repo)
		case err := <-error_channel:
			return nil, err
		case <-quit_channel:
			return repositories, nil
		}
	}
}

func GetCommitsFromTimeRange(repoPath string, since time.Time, until time.Time) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}

	commitMessages := ""

	headRef, err := repo.Head()
	if err != nil {
		return "", err
	}

	commitIter, err := repo.Log(&git.LogOptions{
		From: headRef.Hash(),
	})
	if err != nil {
		return "", err
	}

	err = commitIter.ForEach(func(commit *object.Commit) error {
		if commit.Committer.When.After(since) && commit.Committer.When.Before(until) {
			timeSinceCommit := time.Since(commit.Committer.When)
			formattedTimeSinceCommit := utils.HumanizeDuration(timeSinceCommit)

			commitMessages += fmt.Sprintf("[yellow]%s[white] - %s (%s)\n", commit.Hash.String()[:7], strings.TrimSuffix(commit.Message, "\n"), formattedTimeSinceCommit)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return commitMessages, nil
}
