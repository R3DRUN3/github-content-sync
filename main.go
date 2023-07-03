package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

func main() {

	err := checkEnvVariables()
	if err != nil {
		fmt.Println(err)
		// Handle the error, e.g., exit the program or take appropriate action.
		return
	}

	// All required environment variables are present, continue with your program logic.
	fmt.Println("All environment variables are present.")

	// Read args from env vars
	repoURL := os.Getenv("REPO_URL")
	folder1 := os.Getenv("REPO_FOLDER_1")
	folder2 := os.Getenv("REPO_FOLDER_2")
	token := os.Getenv("GITHUB_TOKEN")

	client := createGitHubClient(token)

	// Check for files present in folder1 but not in folder2
	fmt.Println("Files present in", folder1, "but not in", folder2)
	diffFiles, err := compareFolders(client, repoURL, folder1, folder2)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range diffFiles {
		fmt.Println(*file.Name)
	}

	// Check for files present in both folder1 and folder2
	fmt.Println("\n\nFiles present in both", folder1, "and", folder2, "with newer commits in", folder1)
	newerFiles, err := getFilesWithNewerCommit(client, repoURL, folder1, folder2)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range newerFiles {
		fmt.Println(*file.Name)
	}
}

func checkEnvVariables() error {
	requiredEnvVars := []string{"REPO_URL", "REPO_FOLDER_1", "REPO_FOLDER_2", "GITHUB_TOKEN"}

	for _, envVar := range requiredEnvVars {
		if value, exists := os.LookupEnv(envVar); !exists || value == "" {
			return fmt.Errorf("missing environment variable: %s", envVar)
		}
	}

	return nil
}

func createGitHubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func compareFolders(client *github.Client, repoURL, folder1, folder2 string) ([]*github.RepositoryContent, error) {
	// Extract owner and repo name from the repo URL
	owner, repo := parseRepoURL(repoURL)

	// List files in folder1
	_, folder1Files, _, err := client.Repositories.GetContents(context.Background(), owner, repo, folder1, nil)
	if err != nil {
		return nil, err
	}

	// List files in folder2
	_, folder2Files, _, err := client.Repositories.GetContents(context.Background(), owner, repo, folder2, nil)
	if err != nil {
		return nil, err
	}

	// Find files present in folder1 but not in folder2
	diffFiles := make([]*github.RepositoryContent, 0)
	for _, file1 := range folder1Files {
		found := false
		for _, file2 := range folder2Files {
			if *file1.Name == *file2.Name {
				found = true
				break
			}
		}
		if !found {
			diffFiles = append(diffFiles, file1)
		}
	}

	return diffFiles, nil
}

func getFilesWithNewerCommit(client *github.Client, repoURL, folder1, folder2 string) ([]*github.RepositoryContent, error) {
	// Extract owner and repo name from the repo URL
	owner, repo := parseRepoURL(repoURL)

	// List files in folder1
	_, folder1Files, _, err := client.Repositories.GetContents(context.Background(), owner, repo, folder1, nil)
	if err != nil {
		return nil, err
	}

	// List files in folder2
	_, folder2Files, _, err := client.Repositories.GetContents(context.Background(), owner, repo, folder2, nil)
	if err != nil {
		return nil, err
	}

	// Find files present in both folder1 and folder2 with newer commits in folder1
	newerFiles := make([]*github.RepositoryContent, 0)
	for _, file1 := range folder1Files {
		for _, file2 := range folder2Files {
			if *file1.Name == *file2.Name {
				// Check if the file in folder1 has a newer commit than the file in folder2
				commit1, err := getFileLastCommit(client, owner, repo, folder1, *file1.Name)
				if err != nil {
					return nil, err
				}
				commit2, err := getFileLastCommit(client, owner, repo, folder2, *file2.Name)
				if err != nil {
					return nil, err
				}
				if commit1 != nil && commit2 != nil && commit1.GetCommit().GetCommitter().GetDate().Time.After(commit2.GetCommit().GetCommitter().GetDate().Time) {
					newerFiles = append(newerFiles, file1)
				}
				break
			}
		}
	}

	return newerFiles, nil
}

func getFileLastCommit(client *github.Client, owner, repo, path, file string) (*github.RepositoryCommit, error) {
	commits, _, err := client.Repositories.ListCommits(context.Background(), owner, repo, &github.CommitsListOptions{Path: path + "/" + file})
	if err != nil {
		return nil, err
	}
	if len(commits) > 0 {
		return commits[0], nil
	}
	return nil, nil
}

func parseRepoURL(repoURL string) (string, string) {
	parts := strings.Split(repoURL, "/")
	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]
	repo = strings.TrimSuffix(repo, ".git")
	return owner, repo
}
