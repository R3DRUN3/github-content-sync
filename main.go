package main

import (
	// Standard library packages
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	// Third-party packages
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

func main() {
	err := checkEnvVariables()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("All environment variables are present.")

	// Read environment variables
	repoURL := os.Getenv("REPO_URL")
	folder1 := os.Getenv("REPO_FOLDER_1")
	folder2 := os.Getenv("REPO_FOLDER_2")
	token := os.Getenv("GITHUB_TOKEN")

	// Create a GitHub client with the provided token
	client := createGitHubClient(token)

	fmt.Println("\n[Files present in", folder1, "but not in", folder2, "====>]")
	// Compare folders and get files present in folder1 but not in folder2
	diffFiles, err := compareFolders(client, repoURL, folder1, folder2)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range diffFiles {
		fmt.Println(*file.Name)
	}

	fmt.Println("\n\n[Files present in both", folder1, "and", folder2, "with newer commits in", folder1, "====>]")
	// Get files present in both folder1 and folder2 with newer commits in folder1
	newerFiles, err := getFilesWithNewerCommit(client, repoURL, folder1, folder2)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range newerFiles {
		fmt.Println(*file.Name)
	}
}

// Check if all required environment variables are set
func checkEnvVariables() error {
	requiredEnvVars := []string{"REPO_URL", "REPO_FOLDER_1", "REPO_FOLDER_2", "GITHUB_TOKEN"}

	for _, envVar := range requiredEnvVars {
		if value, exists := os.LookupEnv(envVar); !exists || value == "" {
			return fmt.Errorf("missing environment variable: %s", envVar)
		}
	}

	return nil
}

// Create a GitHub client using the provided token
func createGitHubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// Compare folders and get files present in folder1 but not in folder2
func compareFolders(client *github.Client, repoURL, folder1, folder2 string) ([]*github.RepositoryContent, error) {
	owner, repo := parseRepoURL(repoURL)

	// Get contents of folder1 and folder2 from the GitHub repository
	_, folder1Files, _, err := client.Repositories.GetContents(context.Background(), owner, repo, folder1, nil)
	if err != nil {
		return nil, err
	}

	_, folder2Files, _, err := client.Repositories.GetContents(context.Background(), owner, repo, folder2, nil)
	if err != nil {
		return nil, err
	}

	diffFiles := make([]*github.RepositoryContent, 0)
	diffChan := make(chan *github.RepositoryContent)

	var wg sync.WaitGroup
	wg.Add(len(folder1Files))

	// Compare files in folder1 with files in folder2 concurrently using goroutines
	for _, file1 := range folder1Files {
		go func(file *github.RepositoryContent) {
			defer wg.Done()

			found := false
			for _, file2 := range folder2Files {
				if *file.Name == *file2.Name {
					found = true
					break
				}
			}
			if !found {
				diffChan <- file
			}
		}(file1)
	}

	// Wait for all goroutines to finish and close the diffChan channel
	go func() {
		wg.Wait()
		close(diffChan)
	}()

	// Collect the different files from the diffChan channel
	for file := range diffChan {
		diffFiles = append(diffFiles, file)
	}

	return diffFiles, nil
}

// Get files present in both folder1 and folder2 with newer commits in folder1
func getFilesWithNewerCommit(client *github.Client, repoURL, folder1, folder2 string) ([]*github.RepositoryContent, error) {
	owner, repo := parseRepoURL(repoURL)

	// Get contents of folder1 and folder2 from the GitHub repository
	_, folder1Files, _, err := client.Repositories.GetContents(context.Background(), owner, repo, folder1, nil)
	if err != nil {
		return nil, err
	}

	_, folder2Files, _, err := client.Repositories.GetContents(context.Background(), owner, repo, folder2, nil)
	if err != nil {
		return nil, err
	}

	newerFiles := make([]*github.RepositoryContent, 0)
	newerChan := make(chan *github.RepositoryContent)

	var wg sync.WaitGroup
	wg.Add(len(folder1Files))

	// Compare files in folder1 with files in folder2 and check for newer commits in folder1 using goroutines
	for _, file1 := range folder1Files {
		go func(file *github.RepositoryContent) {
			defer wg.Done()

			for _, file2 := range folder2Files {
				if *file.Name == *file2.Name {
					commit1, err := getFileLastCommit(client, owner, repo, folder1, *file.Name)
					if err != nil {
						log.Println(err)
						return
					}
					commit2, err := getFileLastCommit(client, owner, repo, folder2, *file2.Name)
					if err != nil {
						log.Println(err)
						return
					}
					if commit1 != nil && commit2 != nil && commit1.GetCommit().GetCommitter().GetDate().Time.After(commit2.GetCommit().GetCommitter().GetDate().Time) {
						newerChan <- file
					}
					break
				}
			}
		}(file1)
	}

	// Wait for all goroutines to finish and close the newerChan channel
	go func() {
		wg.Wait()
		close(newerChan)
	}()

	// Collect the files with newer commits from the newerChan channel
	for file := range newerChan {
		newerFiles = append(newerFiles, file)
	}

	return newerFiles, nil
}

// Get the last commit of a file in a specific path
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

// Parse the repository URL and extract the owner and repository name
func parseRepoURL(repoURL string) (string, string) {
	parts := strings.Split(repoURL, "/")
	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]
	repo = strings.TrimSuffix(repo, ".git")
	return owner, repo
}
