FROM mcr.microsoft.com/vscode/devcontainers/go
RUN mkdir /build
ADD ./cmd/main/ /build/cmd/main/
ADD ./cmd/elasticsearch/ /build/cmd/elasticsearch/
ADD ./cmd/filterlines/ /build/cmd/filterlines/
ADD ./cmd/parsecontents/ /build/cmd/parsecontents/
ADD ./cmd/parsedates/ /build/cmd/parsedates/
WORKDIR /build
COPY go.mod go.sum ./
WORKDIR /build/cmd/main
RUN GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o server .
