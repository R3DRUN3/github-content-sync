# GITHUB CONTENT SYNC 🔎 📁
[![goaction](https://github.com/R3DRUN3/github-content-sync/actions/workflows/goaction.yaml/badge.svg)](https://github.com/R3DRUN3/github-content-sync/actions/workflows/goaction.yaml)
[![goreleaser](https://github.com/R3DRUN3/github-content-sync/actions/workflows/release.yaml/badge.svg)](https://github.com/R3DRUN3/github-content-sync/actions/workflows/release.yaml)
[![oci](https://github.com/R3DRUN3/github-content-sync/actions/workflows/oci.yaml/badge.svg)](https://github.com/R3DRUN3/github-content-sync/actions/workflows/oci.yaml)
[![Latest Release](https://img.shields.io/github/release/R3DRUN3/github-content-sync.svg)](https://github.com/R3DRUN3/github-content-sync/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/r3drun3/github-content-sync)](https://goreportcard.com/report/github.com/r3drun3/github-content-sync)  

The *Github Content Sync* tool is a command-line script written in *Go* that allows you to compare the contents of *two folders* within a GitHub repository.  
It helps identify files difference between the two folders.  

<br/>

Basically, if `A` and `B` are the two folders, the tool will output:  
- files present in `A` but not in `B`
- files present in `B` but not in `A`
- files present in both `A` and `B` but with newer commits in `A`  


## Purpose

The purpose of this tool is to facilitate the comparison of folder contents within a GitHub repository.  
**This was specifically meant for those repo that contain documentation in various languages** (divided into different folders) and you need a fast way to know the deltas:  
In this case, usually the reference folder and "*source of truth*" is the "*english*" one (for a real world example take a look at [this repo](https://github.com/cncf/glossary/tree/main/content), for a test playground we use [this one](https://github.com/R3DRUN3/content-sync-tester)).  
Generally, it can be useful in scenarios where you have two folders within a repository and you want to identify the differences between them, such as missing files or files with newer commits.  
## Arguments

The script requires the following environment variables to be set:
- `REPO_URL`: The URL of the GitHub repository to analyze. [MANDATORY]
- `REPO_FOLDER_1`: The name of the reference folder (source of truth, or folder `A`). [MANDATORY]
- `REPO_FOLDER_2`: The name of the second folder to compare to the reference folder (folder `B`). [MANDATORY]
- `TOKEN`: An access token with appropriate permissions to *read* and *open issues* on the target repo. [MANDATORY]
- `OPEN_ISSUE`: If set to `true`, this specify that the script needs to open a "*synchronization issue*" on the target repo, specifying the folder differences. [OPTIONAL]  
The opened issues are structured like [this one](https://github.com/R3DRUN3/content-sync-tester/issues/29).
- `MULTIPLE_ISSUES`: If `OPEN_ISSUE` is set to `true` and this var is also set to `true`, the script will create multiple issues, one for every file difference. [OPTIONAL]  
> **Warning**
> Be careful when setting the `MULTIPLE_ISSUES` var to *true*: if you execute this script against two folders with many files, it will create many issues on your target repo.  


## How it works

The script performs the following steps:
1. Checks the presence of the required environment variables and their values.
1. Creates a GitHub client using the provided access token.
1. Retrieve the content of the two specified folders via the Github client object.
1. Compares the contents of the two specified folders within the repository.
1. Prints the files that are present in the first folder but not in the second folder.
1. Prints the files with newer commits in the first folder compared to the same files in the second folder.
1. Prints the files that are present in the second folder but not in the first folder.
2. If `OPEN_ISSUE` env var is present and set to `true`, opens a "synchronization issue" on the target repo.  
## Examples

You can run this utility in many ways:  

### As an Executable
Download the [release](https://github.com/R3DRUN3/github-content-sync/releases/) that you want and run it:

```shell

export REPO_URL=https://github.com/R3DRUN3/content-sync-tester
export REPO_FOLDER_1=en
export REPO_FOLDER_2=it
export TOKEN=<your-github-token-here>

./github-content-sync
```


Output:
```console
   __   __ _____   _ __  _ __   ___      __   _    _  __ _____   ___   _  __ _____      ___  _  __   _  __   __
 ,'_/  / //_  _/  /// / /// /  / o.)   ,'_/ ,' \  / |/ //_  _/  / _/  / |/ //_  _/    ,' _/ | |/,'  / |/ / ,'_/
/ /_n / /  / /   / ` / / U /  / o \   / /_ / o | / || /  / /   / _/  / || /  / /     _\ `.  | ,'   / || / / /_
|__,'/_/  /_/   /_n_/  \_,'  /___,'   |__/ |_,' /_/|_/  /_/   /___/ /_/|_/  /_/     /___,' /_/    /_/|_/  |__/

[ ALL ENVIRONMENT VARIABLES ARE CONFIGURED ]
[ TARGET REPO URL:  https://github.com/R3DRUN3/content-sync-tester ]

[ FILES PRESENT IN en BUT NOT IN it ]
not_present_in_it.md
not_present_in_it_2.md
test.md


[ FILES PRESENT IN BOTH en AND it WITH NEWER COMMITS IN en ]
doc2.md
last.md


 ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___ ___
/__//__//__//__//__//__//__//__//__//__//__//__//__//__//__//__//__//__//__//__//__//__//__//__//__//__//__//__/
```  

### With Docker (Local Build)
This repo also contain a Dockerfile so you can launch the script as a docker container.  
Clone the repo locally and buil the image:  
```console
git clone https://github.com/r3drun3/github-content-sync \
&& cd github-content-sync \
&& docker build -t github-content-sync:latest .
```  

Run the docker container (change env vars accordingly):  
```console
docker run -it --rm -e REPO_URL=https://github.com/cncf/glossary -e REPO_FOLDER_1=content/en -e REPO_FOLDER_2=content/it -e TOKEN=<your-github-token-here> github-content-sync:latest
```  


### With Docker (Github Packages)
Alternatively, this repo already contains an action to publish the script's OCI image to [Github Packages](https://github.com/features/packages).  
Pull the version that you want: 
```console
docker pull ghcr.io/r3drun3/github-content-sync:1.4.0 
```  

Run the docker container (change env vars accordingly):  
```console
docker run -it --rm -e REPO_URL=https://github.com/cncf/glossary -e REPO_FOLDER_1=content/en -e REPO_FOLDER_2=content/it -e TOKEN=<your-github-token-here> ghcr.io/r3drun3/github-content-sync:1.4.0
```  

### Run via Github Action
The script in this repo can also executed inside a  *Github action*, for an example take a look at the [goaction](https://github.com/R3DRUN3/github-content-sync/actions/workflows/goaction.yaml)  Github Action associated to this repo.  


## Development and Debug
For development and debug I suggest the use of the [VS Code](https://code.visualstudio.com/) IDE.  
In order to debug the script locally, you can create the `.vscode/launch.json` file with the following structure:  
```json
{
    "version": "0.2.0",
    "configurations": [
      {
        "name": "Launch",
        "type": "go",
        "request": "launch",
        "mode": "auto",
        "program": "${workspaceFolder}/main.go",
        "env": {
            "REPO_URL": "<your-github-repo-target-url>",
            "REPO_FOLDER_1": "<path-of-the-reference-folder-inside-target-repo>",
            "REPO_FOLDER_2": "<path-of-the-folder-to-compare-to-the-reference>",
            "TOKEN": "<your-github-token-here>",
            "OPEN_ISSUE": "false",
            "MULTIPLE_ISSUES": "false"
        }
      }
    ]
  }
```  



## Improvements and Next Steps

- It can be useful to maybe add the possibility of comparing multiple folders at the same time, not just 2.


## License

This script is released under the [MIT License](https://opensource.org/license/mit/).  
Feel free to modify and distribute it as per your needs.  


