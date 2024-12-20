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

const ENVIRONMENT_VARIABLE_ERROR = "TINYCARE_WORKSPACE environment variable not set!"

func GetCommits(path string) (string, string, error) {
	if path == "" {
		return ENVIRONMENT_VARIABLE_ERROR, ENVIRONMENT_VARIABLE_ERROR, nil
	}
	repositories, err := FindGitRepositories(path)
	if err != nil {
		return "", "", fmt.Errorf("unable to find git repos for %s: %w", path, err)
	}

	if len(repositories) == 0 {
		return "No Repos Found", "No Repos Found", nil
	}

	dayResult := ""
	weekResult := ""
	for _, repo := range repositories {
		dayCommitMessages, weekCommitMessages, err := GetCommitsFromTimeRange(repo)
		if err != nil {
			return "", "", fmt.Errorf("error pulling commits from repo %s: %w", repo, err)
		}
		if dayCommitMessages != "" {
			dayResult += fmt.Sprintf("[red]%s[white]\n", repo) + dayCommitMessages + "\n"
		}
		if weekCommitMessages != "" {
			weekResult += fmt.Sprintf("[red]%s[white]\n", repo) + weekCommitMessages + "\n"
		}

	}

	return dayResult, weekResult, nil
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

func GetCommitsFromTimeRange(repoPath string) (string, string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", "", err
	}

	headRef, err := repo.Head()
	if err != nil {
		return "", "", err
	}

	commitIter, err := repo.Log(&git.LogOptions{
		From: headRef.Hash(),
	})
	if err != nil {
		return "", "", err
	}
	dayLookBackTime := time.Now().AddDate(0, 0, -1)
	weekLookBackTime := time.Now().AddDate(0, 0, -7)

	dayCommitMessages := ""
	weekCommitMessages := ""

	err = commitIter.ForEach(func(commit *object.Commit) error {
		if commit.Committer.When.After(dayLookBackTime) && commit.Committer.When.Before(time.Now()) {
			timeSinceCommit := time.Since(commit.Committer.When)
			formattedTimeSinceCommit := utils.HumanizeDuration(timeSinceCommit)
			dayCommitMessages += fmt.Sprintf("[yellow]%s[white] - %s (%s)\n", commit.Hash.String()[:7], strings.TrimSuffix(commit.Message, "\n"), formattedTimeSinceCommit)
			weekCommitMessages += fmt.Sprintf("[yellow]%s[white] - %s (%s)\n", commit.Hash.String()[:7], strings.TrimSuffix(commit.Message, "\n"), formattedTimeSinceCommit)
		} else if commit.Committer.When.After(weekLookBackTime) && commit.Committer.When.Before(time.Now()) {
			timeSinceCommit := time.Since(commit.Committer.When)
			formattedTimeSinceCommit := utils.HumanizeDuration(timeSinceCommit)
			weekCommitMessages += fmt.Sprintf("[yellow]%s[white] - %s (%s)\n", commit.Hash.String()[:7], strings.TrimSuffix(commit.Message, "\n"), formattedTimeSinceCommit)
		}
		return nil
	})

	if err != nil {
		return "", "", err
	}

	return dayCommitMessages, weekCommitMessages, nil
}
