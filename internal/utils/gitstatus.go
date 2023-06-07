package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
    "strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetDailyCommits(path string) (string, error) {
	repositories, err := findGitRepositories(path)
	if err != nil {
		return "", err
	}

	result := ""
	for _, repo := range repositories {
		commitMessages, err := getCommitsFromTimeRange(repo, time.Now().AddDate(0, 0, -1), time.Now())
		if err != nil {
			return "", err
		}

		if len(commitMessages) > 0 {
			result += fmt.Sprintf("%s\n", repo)
			result += commitMessages + "\n"
		}
	}

	return result, nil
}

func GetWeeklyCommits(path string) (string, error) {
	repositories, err := findGitRepositories(path)
	if err != nil {
		return "", err
	}

	result := ""
	for _, repo := range repositories {
		commitMessages, err := getCommitsFromTimeRange(repo, time.Now().AddDate(0, 0, -7), time.Now())
		if err != nil {
			return "", err
		}

		if len(commitMessages) > 0 {
			result += fmt.Sprintf("%s\n", repo)
			result += commitMessages + "\n"
		}
	}

	return result, nil
}

func findGitRepositories(path string) ([]string, error) {
	var repositories []string

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			repositories = append(repositories, filepath.Dir(p))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return repositories, nil
}

func getCommitsFromTimeRange(repoPath string, since time.Time, until time.Time) (string, error) {
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
			formattedTimeSinceCommit := humanizeDuration(timeSinceCommit)

			commitMessages += fmt.Sprintf("%s - %s (%s)\n", commit.Hash.String()[:7], strings.TrimSuffix(commit.Message, "\n"), formattedTimeSinceCommit)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return commitMessages, nil
}

func humanizeDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	return fmt.Sprintf("%d hours ago", hours)
}
