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

func GetCommits(path string, lookback int) (string, error) {
	if path == "" {
		return "TINYCARE_WORKSPACE environment variable not set!", nil
	}
	repositories, err := FindGitRepositories(path)
	if err != nil {
		return "", fmt.Errorf("unable to find git repos for %s: %w", path, err)
	}

	if len(repositories) == 0 {
		return "No Repos Found", nil
	}

	result := ""
	for _, repo := range repositories {
		commitMessages, err := GetCommitsFromTimeRange(repo, time.Now().AddDate(0, 0, lookback), time.Now())
		if err != nil {
			return "", fmt.Errorf("error pulling commits from repo %s: %w", repo, err)
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
			defer wg.Done()
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
		}(path)
	}
	wg.Wait()
	close(q)
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
			close(repo_channel)
			close(error_channel)
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
