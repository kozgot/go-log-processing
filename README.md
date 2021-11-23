![Go](https://github.com/kozgot/go-log-processing/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/kozgot/go-log-processing)](https://goreportcard.com/report/github.com/kozgot/go-log-processing)

# Diagnostic log processing and analysis

This repository contains the source code of my thesis project, which is Diagnostic log data processing

## Development
* For development, the VSCode IDE is used with the Go and the Remote-Containers extension.
* Download the Remote-Containers extension here: [Remote-Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
   - The Remote - Containers extension lets you use a Docker container as a full-featured development environment.
   - `Tip: Allocate more than 2GB memory for docker, for debbuging inside the Remote container 2 GB won't be enough and the attach is going to keep disconnecting. `
   - In VSCode press `F1`, then select ```Remote-Containers: Open Folder in Container``` and select the folder of one of the microservices (eg.: parser).
   - After the container starts, open a new VSCode window, press `F1`, then select ```Remote-Containers: Open Folder in Container``` and select the folder of another microservice (eg.: postprocessor). 
   - Finally, do the same for elasticuploader (new window, then open in container...)
   - Press `F5` in all three windows to start a debug session. You need to start the services in a specific order: elasticuploader first, then postprocessor, then parser.
* Docker Desktop and Docker Compose are also needed for development

