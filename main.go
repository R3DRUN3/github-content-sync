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
	folder1Branch := os.Getenv("FOLDER_1_BRANCH")
	if folder1Branch == "" {
		folder1Branch = "main"
	}
	folder2Branch := os.Getenv("FOLDER_2_BRANCH")
	if folder2Branch == "" {
		folder2Branch = "main"
	}

	// Create a GitHub client with the provided token
	client := createGitHubClient(token)
	fmt.Println("[ TARGET REPO URL: ", repoURL, "]")
	fmt.Println("\n[ FILES PRESENT IN", folder1, "ON BRANCH", folder1Branch, "BUT NOT IN", folder2, "ON BRANCH", folder2Branch, "]")
	// Compare folders and get files present in folder1 but not in folder2
	diffFiles, newerFiles, diffFilesFolder2, err := compareFolders(client, repoURL, folder1, folder1Branch, folder2, folder2Branch)
	if err != nil {
		log.Fatal(err)
	}
	printFilesSorted(diffFiles)

	fmt.Println("\n\n[ FILES PRESENT IN BOTH", folder1, "AND", folder2, "WITH NEWER COMMITS IN", folder1, "]")
	// Print files present in both folder1 and folder2 with newer commits in folder1
	printFilesSorted(newerFiles)

	fmt.Println("\n\n[ FILES PRESENT IN", folder2, "ON BRANCH", folder2Branch, "BUT NOT IN", folder1, "ON BRANCH", folder1Branch, "]")
	// Print files present in folder2 but not in folder1
	printFilesSorted(diffFilesFolder2)

	// Open an issue if OPEN_ISSUE env var is set to true
	if os.Getenv("OPEN_ISSUE") == "true" {
		err := openSyncIssue(client, repoURL, folder1, folder2, diffFiles, diffFilesFolder2, newerFiles)
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
	requiredEnvVars := []string{"REPO_URL", "REPO_FOLDER_1", "REPO_FOLDER_2", "TOKEN"}
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
// Compare folders and get files present in folder1 but not in folder2,
// files with newer commits in folder1, and files present in folder2 but not in folder1
func compareFolders(client *github.Client, repoURL, folder1, folder1Branch, folder2, folder2Branch string) ([]*github.RepositoryContent, []*github.RepositoryContent, []*github.RepositoryContent, error) {
	owner, repo := parseRepoURL(repoURL)

	// Get contents of folder1 and folder2 from the GitHub repository
	folder1Files, err := getFolderContents(client, owner, repo, folder1, folder1Branch)
	if err != nil {
		return nil, nil, nil, err
	}

	folder2Files, err := getFolderContents(client, owner, repo, folder2, folder2Branch)
	if err != nil {
		return nil, nil, nil, err
	}

	diffFilesFolder1 := make([]*github.RepositoryContent, 0)
	newerFiles := make([]*github.RepositoryContent, 0)
	diffFilesFolder2 := make([]*github.RepositoryContent, 0)

	var wg sync.WaitGroup
	wg.Add(len(folder1Files))

	var mu sync.Mutex // Declare a mutex

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
				mu.Lock() // Lock the mutex before modifying the slice
				diffFilesFolder1 = append(diffFilesFolder1, file)
				mu.Unlock() // Unlock the mutex after modifying the slice
			} else {
				commit1, err := getFileLastCommit(client, owner, repo, folder1, folder1Branch, *file.Name)
				if err != nil {
					log.Println(err)
					return
				}
				commit2, err := getFileLastCommit(client, owner, repo, folder2, folder2Branch, *file.Name)
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

	// Compare files in folder2 with files in folder1
	for _, file2 := range folder2Files {
		found := false
		for _, file1 := range folder1Files {
			if *file2.Name == *file1.Name {
				found = true
				break
			}
		}
		if !found {
			diffFilesFolder2 = append(diffFilesFolder2, file2)
		}
	}

	wg.Wait()

	return diffFilesFolder1, newerFiles, diffFilesFolder2, nil
}

// Get contents of a folder from the GitHub repository
func getFolderContents(client *github.Client, owner, repo, folder, branch string) ([]*github.RepositoryContent, error) {
	opt := &github.RepositoryContentGetOptions{
		Ref: branch,
	}
	_, files, _, err := client.Repositories.GetContents(context.Background(), owner, repo, folder, opt)
	if err != nil {
		return nil, err
	}
	// Check if recursive scanning is enabled
	if recursive, _ := strconv.ParseBool(os.Getenv("RECURSIVE")); recursive {
		var allFiles []*github.RepositoryContent
		for _, file := range files {
			if *file.Type == "dir" {
				subFiles, err := getFolderContentsRecursive(client, owner, repo, *file.Path, branch)
				if err != nil {
					return nil, err
				}
				allFiles = append(allFiles, subFiles...)
			} else {
				allFiles = append(allFiles, file)
			}
		}
		return allFiles, nil
	}
	return files, nil
}

// Recursive helper function to get contents of a directory
func getFolderContentsRecursive(client *github.Client, owner, repo, path, branch string) ([]*github.RepositoryContent, error) {
	opt := &github.RepositoryContentGetOptions{
		Ref: branch,
	}
	_, files, _, err := client.Repositories.GetContents(context.Background(), owner, repo, path, opt)
	if err != nil {
		return nil, err
	}

	var allFiles []*github.RepositoryContent
	for _, file := range files {
		if *file.Type == "dir" {
			subFiles, err := getFolderContentsRecursive(client, owner, repo, *file.Path, branch)
			if err != nil {
				return nil, err
			}
			allFiles = append(allFiles, subFiles...)
		} else {
			allFiles = append(allFiles, file)
		}
	}
	return allFiles, nil
}

// Get the last commit of a file in a specific path on a particular branch
func getFileLastCommit(client *github.Client, owner, repo, path, branch, file string) (*github.RepositoryCommit, error) {
	commits, _, err := client.Repositories.ListCommits(context.Background(), owner, repo, &github.CommitsListOptions{Path: path + "/" + file, SHA: branch})
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
func openSyncIssue(client *github.Client, repoURL, folder1, folder2 string, diffFiles, diffFilesFolder2, newerFiles []*github.RepositoryContent) error {
	owner, repo := parseRepoURL(repoURL)
	// Check if  need to create multiple issues
	if os.Getenv("MULTIPLE_ISSUES") == "true" {
		for _, file := range diffFiles {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			issueTitle := "Synchronization Issue [" + timestamp + "]: " + folder1 + " vs " + folder2
			issueBody := "## Synchronization Issue\n\n" +
				"Folder1: " + folder1 + "\n\n" +
				"Folder2: " + folder2 + "\n\n" +
				"### Files present in " + folder1 + " but not in " + folder2 + "\n"
			issueBody += "- " + *file.Name + "\n"
			issueRequest := &github.IssueRequest{
				Title: &issueTitle,
				Body:  &issueBody,
			}
			_, _, err := client.Issues.Create(context.Background(), owner, repo, issueRequest)
			if err != nil {
				return err
			}
			fmt.Println("\n[ SYNCHRONIZATION ISSUE OPENED ]")
		}
		for _, file := range newerFiles {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			issueTitle := "Synchronization Issue [" + timestamp + "]: " + folder1 + " vs " + folder2
			issueBody := "## Synchronization Issue\n\n" +
				"Folder1: " + folder1 + "\n\n" +
				"Folder2: " + folder2 + "\n\n" +
				"### Files present in both " + folder1 + " and " + folder2 + " with newer commits in " + folder1 + "\n"
			issueBody += "- " + *file.Name + "\n"
			issueRequest := &github.IssueRequest{
				Title: &issueTitle,
				Body:  &issueBody,
			}
			_, _, err := client.Issues.Create(context.Background(), owner, repo, issueRequest)
			if err != nil {
				return err
			}
			fmt.Println("\n[ SYNCHRONIZATION ISSUE OPENED ]")
		}
		for _, file := range diffFilesFolder2 {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			issueTitle := "Synchronization Issue [" + timestamp + "]: " + folder1 + " vs " + folder2
			issueBody := "## Synchronization Issue\n\n" +
				"Folder1: " + folder1 + "\n\n" +
				"Folder2: " + folder2 + "\n\n" +
				"### Files present in " + folder2 + " but not in " + folder1 + "\n"
			issueBody += "- " + *file.Name + "\n"
			issueRequest := &github.IssueRequest{
				Title: &issueTitle,
				Body:  &issueBody,
			}
			_, _, err := client.Issues.Create(context.Background(), owner, repo, issueRequest)
			if err != nil {
				return err
			}
			fmt.Println("\n[ SYNCHRONIZATION ISSUE OPENED ]")
		}
	} else { // create a single issue
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

		issueBody += "\n### Files present in " + folder2 + " but not in " + folder1 + "\n"
		for _, file := range diffFilesFolder2 {
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
	return nil
}
