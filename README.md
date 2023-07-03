# github-content-sync 🔎 📁
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/r3drun3/github-content-sync)](https://goreportcard.com/report/github.com/r3drun3/github-content-sync)  

The *Github Content Sync* tool is a command-line script written in *Go* that allows you to compare the contents of two folders in a GitHub repository.  
It helps identify files that are present in one folder but not in another, as well as files that have newer commits in one folder compared to another.  
## Purpose

The purpose of this tool is to facilitate the comparison of folder contents within a GitHub repository, expecially for those repo that contain documentation in various languages (divided into different folders).  
It can be useful in scenarios where you have two folders within a repository and you want to identify the differences between them, such as missing files or files with newer commits.  
## Arguments

The script requires the following environment variables to be set: 
- `REPO_URL`: The URL of the GitHub repository to analyze. 
- `REPO_FOLDER_1`: The name of the first folder to compare. 
- `REPO_FOLDER_2`: The name of the second folder to compare. 
- `GITHUB_TOKEN`: An access token with appropriate permissions to access the repository.
## How it works

The script performs the following steps:
1. Checks the presence of the required environment variables and their values.
2. Creates a GitHub client using the provided access token.
3. Compares the contents of the two specified folders within the repository.
4. Prints the files that are present in the first folder but not in the second folder.
5. Retrieves files that exist in both folders and have newer commits in the first folder.
6. Prints the files with newer commits in the first folder compared to the second folder.
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
All environment variables are present.

[Files present in content/en but not in content/it ====>]
_TEMPLATE.md
blue-green-deployment.md
application-programming-interface.md
auto-scaling.md
bare-metal-machine.md
cloud-native-security.md
chaos-engineering.md
cloud-computing.md
cloud-native-apps.md
container-image.md
container.md
container-orchestration.md
contributor-ladder
continuous-delivery.md
continuous-deployment.md
continuous-integration.md
data-center.md
database-as-a-service.md
distributed-systems.md
digital-certificate.md
edge-computing.md
event-streaming.md
function-as-a-service.md
gitops.md
horizontal-scaling.md
hypervisor.md
kubernetes.md
microservices-architecture.md
load-balancer.md
mutual-transport-layer-security.md
pod.md
policy-as-code.md
serverless.md
role-based-access-control.md
search.md
security-chaos-engineering.md
service-discovery.md
service-proxy.md
transport-layer-security.md
stateless-apps.md
vertical-scaling.md
virtualization.md
zero-trust-architecture.md


[Files present in both content/en and content/it with newer commits in content/en ====>]
agile-software-development.md
distributed-apps.md
infrastructure-as-code.md
devsecops.md
observability.md
self-healing.md
debugging.md
virtual-machine.md
platform-as-a-service.md
client-server-architecture.md
managed-services.md
site-reliability-engineering.md
cluster.md
version-control.md
software-as-a-service.md
scalability.md
loosely-coupled-architecture.md
_index.md
nodes.md
portability.md
monolithic-apps.md
containers-as-a-service.md
service.md
event-driven-architecture.md
multitenancy.md
idempotence.md
style-guide
tightly-coupled-architectures.md
contribute
canary-deployment.md
abstraction.md
reliability.md
firewall.md
```  





## Dependencies

The script uses the following external dependencies: 
- [go-github](https://github.com/google/go-github) : A Go library for accessing the GitHub API.

Please refer to the [Go documentation](https://golang.org/doc/)  for instructions on how to install and manage dependencies.
## License

This script is released under the [MIT License](https://chat.openai.com/LICENSE).  
Feel free to modify and distribute it as per your needs.  


