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
	if len(os.Args) != 4 {
		fmt.Println("Usage: github-content-sync <repo-url> <folder1> <folder2>")
		return
	}

	repoURL := os.Args[1]
	folder1 := os.Args[2]
	folder2 := os.Args[3]
	token := os.Getenv("GITHUB_TOKEN") // Read Github token form env variable

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

	// Check for files in folder1 with newer last update than folder2
	// fmt.Println("\nFiles in", folder1, "with newer last update than", folder2)
	// updatedFiles, err := compareLastUpdate(client, repoURL, folder1, folder2)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, file := range updatedFiles {
	// 	fmt.Println(*file.Name)
	// }

	// // Parallelize the execution of comparing last update
	// var wg sync.WaitGroup
	// wg.Add(2)
	// go func() {
	// 	defer wg.Done()
	// 	// Compare folder1 -> folder2
	// 	compareLastUpdate(client, repoURL, folder1, folder2)
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	// Compare folder2 -> folder1
	// 	compareLastUpdate(client, repoURL, folder2, folder1)
	// }()
	// wg.Wait()
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

// func compareLastUpdate(client *github.Client, repoURL, folder1, folder2 string) ([]*github.RepositoryContent, error) {
// 	// Extract owner and repo name from the repo URL
// 	owner, repo := parseRepoURL(repoURL)

// 	// List commits in the repository
// 	commits, _, err := client.Repositories.ListCommits(context.Background(), owner, repo, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Retrieve the commit history for each file in folder1
// 	commitHistory := make(map[string]time.Time)
// 	for _, commit := range commits {
// 		files, _, _, err := client.Repositories.ListFiles(context.Background(), owner, repo, *commit.SHA, nil)
// 		if err != nil {
// 			return nil, err
// 		}
// 		for _, file := range files {
// 			if strings.HasPrefix(*file.Path, folder1) {
// 				if lastUpdate, ok := commitHistory[*file.Path]; !ok || lastUpdate.Before(*commit.Commit.Author.Date) {
// 					commitHistory[*file.Path] = *commit.Commit.Author.Date
// 				}
// 			}
// 		}
// 	}

// 	// List files in folder2
// 	_, folder2Files, _, err := client.Repositories.GetContents(context.Background(), owner, repo, folder2, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Find files in folder1 with newer last update than folder2
// 	updatedFiles := make([]*github.RepositoryContent, 0)
// 	for _, file := range folder2Files {
// 		if lastUpdate, ok := commitHistory[*file.Path]; ok {
// 			if lastUpdate.After(file.GetCommit().GetCommit().Author.GetDate().Time) {
// 				updatedFiles = append(updatedFiles, file)
// 			}
// 		}
// 	}

// 	return updatedFiles, nil
// }

func parseRepoURL(repoURL string) (string, string) {
	parts := strings.Split(repoURL, "/")
	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]
	repo = strings.TrimSuffix(repo, ".git")
	return owner, repo
}
