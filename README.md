# github-content-sync 🔎 📁
[![lint](https://github.com/R3DRUN3/github-content-sync/actions/workflows/lint.yaml/badge.svg)](https://github.com/R3DRUN3/github-content-sync/actions/workflows/lint.yaml)
[![goreleaser](https://github.com/R3DRUN3/github-content-sync/actions/workflows/release.yaml/badge.svg)](https://github.com/R3DRUN3/github-content-sync/actions/workflows/release.yaml)
[![oci](https://github.com/R3DRUN3/github-content-sync/actions/workflows/oci.yaml/badge.svg)](https://github.com/R3DRUN3/github-content-sync/actions/workflows/oci.yaml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/r3drun3/github-content-sync)](https://goreportcard.com/report/github.com/r3drun3/github-content-sync)  

The *Github Content Sync* tool is a command-line script written in *Go* that allows you to compare the contents of two folders within a GitHub repository.  
It helps identify files that are present in one folder but not in another, as well as files that have newer commits in one folder compared to another.  
## Purpose

The purpose of this tool is to facilitate the comparison of folder contents within a GitHub repository.  
This was specifically meant for those repo that contain documentation in various languages (divided into different folders) and you need a fast way to know the deltas.  
Generally, it can be useful in scenarios where you have two folders within a repository and you want to identify the differences between them, such as missing files or files with newer commits.  
## Arguments

The script requires the following environment variables to be set: 
- `REPO_URL`: The URL of the GitHub repository to analyze. 
- `REPO_FOLDER_1`: The name of the first folder to compare. 
- `REPO_FOLDER_2`: The name of the second folder to compare. 
- `GITHUB_TOKEN`: An access token with appropriate permissions to access the repository.
## How it works

The script performs the following steps:
1. Checks the presence of the required environment variables and their values.
1. Creates a GitHub client using the provided access token.
1. Retrieve the content of the two specified folders via the Github client object.
1. Compares the contents of the two specified folders within the repository.
1. Retrieves files that exist in both folders and have newer commits in the first folder.
1. Prints the files that are present in the first folder but not in the second folder.
1. Prints the files with newer commits in the first folder compared to the same files in the second folder.
## Examples

Here are some examples of how to use the Folder Comparison Tool:
1. Compare two folders within a GitHub repository:

```shell

export REPO_URL=https://github.com/cncf/glossary
export REPO_FOLDER_1=content/en
export REPO_FOLDER_2=content/it
export GITHUB_TOKEN=your-github-token

go run main.go
```


Output:
```console
#################### [GITHUB CONTENT SYNC] ####################
All environment variables are present.

[Files present in content/en but not in content/it ====>]
_TEMPLATE.md
application-programming-interface.md
auto-scaling.md
bare-metal-machine.md
blue-green-deployment.md
chaos-engineering.md
cloud-computing.md
cloud-native-apps.md
cloud-native-security.md
container-image.md
container-orchestration.md
container.md
continuous-delivery.md
continuous-deployment.md
continuous-integration.md
contributor-ladder
data-center.md
database-as-a-service.md
digital-certificate.md
distributed-systems.md
edge-computing.md
event-streaming.md
function-as-a-service.md
gitops.md
horizontal-scaling.md
hypervisor.md
kubernetes.md
load-balancer.md
microservices-architecture.md
mutual-transport-layer-security.md
pod.md
policy-as-code.md
role-based-access-control.md
search.md
security-chaos-engineering.md
serverless.md
service-discovery.md
service-proxy.md
stateless-apps.md
transport-layer-security.md
vertical-scaling.md
virtualization.md
zero-trust-architecture.md


[Files present in both content/en and content/it with newer commits in content/en ====>]
_index.md
abstraction.md
agile-software-development.md
canary-deployment.md
client-server-architecture.md
cluster.md
containers-as-a-service.md
contribute
debugging.md
devsecops.md
distributed-apps.md
event-driven-architecture.md
firewall.md
idempotence.md
infrastructure-as-code.md
loosely-coupled-architecture.md
managed-services.md
monolithic-apps.md
multitenancy.md
nodes.md
observability.md
platform-as-a-service.md
portability.md
reliability.md
scalability.md
self-healing.md
service.md
site-reliability-engineering.md
software-as-a-service.md
style-guide
tightly-coupled-architectures.md
version-control.md
virtual-machine.md

###############################################################
```  

## With Docker (Local Build)
This repo also contain a Dockerfile so you can launch the script as a docker container.  
buil the image:  
```console
docker build -t github-content-sync:v1 .
```  

Run the docker container (change env vars accordingly):  
```console
docker run -it --rm -e REPO_URL=https://github.com/cncf/glossary -e REPO_FOLDER_1=content/en -e REPO_FOLDER_2=content/it -e GITHUB_TOKEN=<your-token-here> github-content-sync:v1
```  





## Dependencies

The script uses the following external dependencies: 
- [go-github](https://github.com/google/go-github) : A Go library for accessing the GitHub API.

Please refer to the [Go documentation](https://golang.org/doc/)  for instructions on how to install and manage dependencies.
## License

This script is released under the [MIT License](https://chat.openai.com/LICENSE).  
Feel free to modify and distribute it as per your needs.  


