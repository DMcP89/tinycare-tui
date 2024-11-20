package local

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func Test_GetCommits(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	defer removeTempDir(t, tempDir)

	// Initialize a git repository in the temporary directory
	repoPath := initGitRepo(t, tempDir)

	// Create a new commit with the current timestamp
	hash := createCommit(t, repoPath, "Test commit", time.Now())

	// Call the GetDailyCommits function
	commits, err := GetCommits(repoPath, -1)
	if err != nil {
		t.Errorf("GetDailyCommits returned an error: %v", err)
	}

	// Verify the commit message
	expectedMessage := fmt.Sprintf("[red]%s[white]\n[yellow]%s[white] - Test commit ([green]0 h(s) ago[white])\n\n", repoPath, hash.String()[:7])
	if commits != expectedMessage {
		t.Errorf("Expected commit message '%s', got '%s'", expectedMessage, commits)
	}
}

// Helper functions for testing

func createTempDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "git_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	return tempDir
}

func removeTempDir(t *testing.T, tempDir string) {
	err := os.RemoveAll(tempDir)
	if err != nil {
		t.Fatalf("Failed to remove temporary directory: %v", err)
	}
}

func initGitRepo(t *testing.T, repoPath string) string {
	_, err := git.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}
	return repoPath
}

func createCommit(t *testing.T, repoPath string, message string, timestamp time.Time) plumbing.Hash {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		t.Fatalf("Failed to open git repository: %v", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Failed to get worktree: %v", err)
	}

	// Create a new file
	filePath := filepath.Join(repoPath, "test.txt")
	err = os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Add the file to the repository
	_, err = wt.Add("test.txt")
	if err != nil {
		t.Fatalf("Failed to add file to repository: %v", err)
	}

	// Commit the changes
	hash, err := wt.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  timestamp,
		},
	})
	if err != nil {
		t.Fatalf("Failed to commit changes: %v", err)
	}
	return hash

}

func Test_FindGitRepositories(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := createTempDir(t)
	defer removeTempDir(t, tempDir)

	// Create some test repositories in the temporary directory
	repoPaths := []string{
		initGitRepo(t, filepath.Join(tempDir, "repo1")),
		initGitRepo(t, filepath.Join(tempDir, "repo2")),
		initGitRepo(t, filepath.Join(tempDir, "repo3")),
	}

	// Call the FindGitRepositories function
	repos, err := FindGitRepositories(tempDir)
	if err != nil {
		t.Fatalf("Failed to find git repositories: %v", err)
	}

	// Assert that the correct number of repositories were found
	if len(repos) != len(repoPaths) {
		t.Errorf("Expected %d repositories, got %d", len(repoPaths), len(repos))
	}

	// Assert that the found repositories match the expected paths
	for _, repoPath := range repoPaths {
		found := false
		for _, repo := range repos {
			if strings.Contains(repo, repoPath) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected repository %s not found", repoPath)
		}
	}
}

func Test_GetCommitsFromTimeRange(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := createTempDir(t)
	defer removeTempDir(t, tempDir)

	// Initialize a git repository in the temporary directory
	repoPath := initGitRepo(t, tempDir)

	// Create some commits with different timestamps
	now := time.Now()
	createCommit(t, repoPath, "Commit 1", now.Add(-7*24*time.Hour))
	createCommit(t, repoPath, "Commit 2", now.Add(-6*24*time.Hour))
	createCommit(t, repoPath, "Commit 3", now.Add(-5*24*time.Hour))
	createCommit(t, repoPath, "Commit 4", now.Add(-4*24*time.Hour))
	createCommit(t, repoPath, "Commit 5", now.Add(-3*24*time.Hour))
	createCommit(t, repoPath, "Commit 6", now.Add(-2*24*time.Hour))
	createCommit(t, repoPath, "Commit 7", now.Add(-24*time.Hour))
	createCommit(t, repoPath, "Commit 8", now)

	// Call the GetCommitsFromTimeRange function
	_, err := GetCommitsFromTimeRange(repoPath, now.Add(-7*24*time.Hour), now)
	if err != nil {
		t.Fatalf("Failed to get commits from time range: %v", err)
	}
}
