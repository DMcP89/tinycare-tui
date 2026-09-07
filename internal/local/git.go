package local

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"log/slog"
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
			err := filepath.WalkDir(path, func(p string, info fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() && info.Name() == ".git" {
					c <- filepath.Dir(p)
					return filepath.SkipDir
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

var (
	repoCacheMu    sync.Mutex
	repoCachePath  string
	repoCacheRepos []string
)

var logger = slog.Default()

func FindGitRepositories(path string) ([]string, error) {
	logger.Info("scanned directory for git repositories", "path", path)
	repoCacheMu.Lock()
	defer repoCacheMu.Unlock()
	if path == repoCachePath && repoCacheRepos != nil {
		return repoCacheRepos, nil
	}

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
			repoCachePath = path
			repoCacheRepos = repositories
			return repositories, nil
		}
	}
}

func GetCommitsFromTimeRange(repoPath string) (string, string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		logger.Error("failed to open repository", "path", repoPath, "error", err)
		return "", "", err
	}

	headRef, err := repo.Head()
	if err != nil {
		logger.Warn("failed to get head of repository", "path", repoPath, "error", err)
		return "", "", err
	}

	dayLookBackTime := time.Now().AddDate(0, 0, -1)
	weekLookBackTime := time.Now().AddDate(0, 0, -7)

	commitIter, err := repo.Log(&git.LogOptions{
		From:  headRef.Hash(),
		Since: &weekLookBackTime,
	})
	if err != nil {
		return "", "", err
	}

	dayCommitMessages := ""
	weekCommitMessages := ""

	err = commitIter.ForEach(func(commit *object.Commit) error {
		if commit.Committer.When.After(dayLookBackTime) {
			timeSinceCommit := time.Since(commit.Committer.When)
			formattedTimeSinceCommit := utils.HumanizeDuration(timeSinceCommit)
			dayCommitMessages += fmt.Sprintf("[yellow]%s[white] - %s (%s)\n", commit.Hash.String()[:7], strings.TrimSuffix(commit.Message, "\n"), formattedTimeSinceCommit)
			weekCommitMessages += fmt.Sprintf("[yellow]%s[white] - %s (%s)\n", commit.Hash.String()[:7], strings.TrimSuffix(commit.Message, "\n"), formattedTimeSinceCommit)
		} else {
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
