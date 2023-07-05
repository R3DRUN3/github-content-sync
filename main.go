package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

func main() {
	header := figure.NewFigure("GITHUB CONTENT SYNC", "eftitalic", true)
	header.Print()
	envVars, err := getEnvVariables()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Read environment variables
	repoURL := envVars[0]
	folder1 := envVars[1]
	folder2 := envVars[2]
	token := envVars[3]

	// Create a GitHub client with the provided token
	client := createGitHubClient(token)
	fmt.Println("[ TARGET REPO URL: ", repoURL, "]")
	fmt.Println("\n[ FILES PRESENT IN", folder1, "BUT NOT IN", folder2, "]")
	// Compare folders and get files present in folder1 but not in folder2
	diffFiles, newerFiles, err := compareFolders(client, repoURL, folder1, folder2)
	if err != nil {
		log.Fatal(err)
	}
	printFilesSorted(diffFiles)

	fmt.Println("\n\n[ FILES PRESENT IN BOTH", folder1, "AND", folder2, "WITH NEWER COMMITS IN", folder1, "]")
	// Print files present in both folder1 and folder2 with newer commits in folder1
	printFilesSorted(newerFiles)

	// Open an issue if OPEN_ISSUE env var is set to true
	if os.Getenv("OPEN_ISSUE") == "true" {
		err := openSyncIssue(client, repoURL, folder1, folder2, diffFiles, newerFiles)
		if err != nil {
			log.Fatal(err)
		}
	}

	footer := figure.NewFigure("----------------------------", "eftitalic", true)
	footer.Print()
	fmt.Println()
}

// Check if all required environment variables are set and return the list of values
func getEnvVariables() ([]string, error) {
	requiredEnvVars := []string{"REPO_URL", "REPO_FOLDER_1", "REPO_FOLDER_2", "GITHUB_TOKEN"}
	envVarValues := make([]string, len(requiredEnvVars))

	for i, envVar := range requiredEnvVars {
		if value, exists := os.LookupEnv(envVar); !exists || value == "" {
			return nil, fmt.Errorf("missing environment variable ===> %s", envVar)
		} else {
			envVarValues[i] = value
		}
	}

	fmt.Println("\n[ ALL ENVIRONMENT VARIABLES ARE CONFIGURED ]")
	return envVarValues, nil
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

// Compare folders and get files present in folder1 but not in folder2, and files with newer commits in folder1
func compareFolders(client *github.Client, repoURL, folder1, folder2 string) ([]*github.RepositoryContent, []*github.RepositoryContent, error) {
	owner, repo := parseRepoURL(repoURL)

	// Get contents of folder1 and folder2 from the GitHub repository
	folder1Files, err := getFolderContents(client, owner, repo, folder1)
	if err != nil {
		return nil, nil, err
	}

	folder2Files, err := getFolderContents(client, owner, repo, folder2)
	if err != nil {
		return nil, nil, err
	}

	diffFiles := make([]*github.RepositoryContent, 0)
	newerFiles := make([]*github.RepositoryContent, 0)

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
				diffFiles = append(diffFiles, file)
			} else {
				commit1, err := getFileLastCommit(client, owner, repo, folder1, *file.Name)
				if err != nil {
					log.Println(err)
					return
				}
				commit2, err := getFileLastCommit(client, owner, repo, folder2, *file.Name)
				if err != nil {
					log.Println(err)
					return
				}
				if commit1 != nil && commit2 != nil && commit1.GetCommit().GetCommitter().GetDate().Time.After(commit2.GetCommit().GetCommitter().GetDate().Time) {
					newerFiles = append(newerFiles, file)
				}
			}
		}(file1)
	}

	wg.Wait()

	return diffFiles, newerFiles, nil
}

// Get contents of a folder from the GitHub repository
func getFolderContents(client *github.Client, owner, repo, folder string) ([]*github.RepositoryContent, error) {
	_, files, _, err := client.Repositories.GetContents(context.Background(), owner, repo, folder, nil)
	if err != nil {
		return nil, err
	}
	return files, nil
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

// Print the files in lexicographic order
func printFilesSorted(files []*github.RepositoryContent) {
	// Sort files by their names
	sort.Slice(files, func(i, j int) bool {
		return *files[i].Name < *files[j].Name
	})

	// Print the sorted file names
	for _, file := range files {
		fmt.Println(*file.Name)
	}
}

// Open a synchronization issue on GitHub repository
// Open a synchronization issue on GitHub repository
func openSyncIssue(client *github.Client, repoURL, folder1, folder2 string, diffFiles, newerFiles []*github.RepositoryContent) error {
	owner, repo := parseRepoURL(repoURL)

	// Generate timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	issueTitle := "Synchronization Issue [" + timestamp + "]: " + folder1 + " vs " + folder2
	issueBody := "## Synchronization Issue\n\n" +
		"Folder1: " + folder1 + "\n\n" +
		"Folder2: " + folder2 + "\n\n" +
		"### Files present in " + folder1 + " but not in " + folder2 + "\n"
	for _, file := range diffFiles {
		issueBody += "- " + *file.Name + "\n"
	}

	issueBody += "\n### Files present in both " + folder1 + " and " + folder2 + " with newer commits in " + folder1 + "\n"
	for _, file := range newerFiles {
		issueBody += "- " + *file.Name + "\n"
	}

	issueRequest := &github.IssueRequest{
		Title: &issueTitle,
		Body:  &issueBody,
	}

	_, _, err := client.Issues.Create(context.Background(), owner, repo, issueRequest)
	if err != nil {
		return err
	}

	fmt.Println("\n[ SYNCHRONIZATION ISSUE OPENED ]")
	return nil
}
