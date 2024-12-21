package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		name   string
		docker bool
		make   bool
	)

	flag.StringVar(&name, "name", "Project Name", "Your name")
	flag.BoolVar(&docker, "docker", false, "Add a Dockerfile")
	flag.BoolVar(&make, "make", false, "Add a Makefile")

	flag.Parse()

	fmt.Printf("Project name: %s\n", name)
	fmt.Printf("Dockerize? %t\n", docker)
	fmt.Printf("Makefile? %t\n", make)

	projectName := name

	if docker {
		createDockerFile(projectName)
	}

	if make {
		createMakefile(projectName)
	}
}

func createDockerFile(projectName string) {

	baseDockerfile := fmt.Sprintf(`
FROM golang:1.23-alpine AS build

WORKDIR /usr/src/app

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/bin/kbgbs ./cmd/main.go

EXPOSE 8333

ENTRYPOINT ["./app/%s"]

	`, projectName)

	e := os.WriteFile("Dockerfile", []byte(baseDockerfile), 0755)
	if e != nil {
		fmt.Printf("Unable to create Dockerfile: %v", e)
	}
}

func createMakefile(projectName string) {

	baseMakefile := fmt.Sprintf(`
DOCKER_IMAGE=%s
DOCKER_NETWORK=%s

new-workstation: docker-pg run
	@echo "workstation ready"

run: build
	@./bin/main

build:
	@go build -o bin/main cmd/main.go

docker-run:
	@docker run --rm --name ${DOCKER_IMAGE} \ 
	--network=kbgbs \
	-p 8333:8333 \
	--env-file .env_docker \
	-d ${DOCKER_IMAGE}

docker-build:
	@docker build \
	-t ${DOCKER_IMAGE} \
	--no-cache .

shell-into:
	@docker run -it \
	--rm --name ${DOCKER_IMAGE} \
	--network=${DOCKER_NETWORK} \
	-v $(PWD):/app  \
	--entrypoint /bin/bash \
	${DOCKER_IMAGE}

docker-pg:
	@docker run \
	--rm --name postgres \
	--network==${DOCKER_NETWORK} \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_DB=postgres \
	-p 5432:5432 \
	-d postgres:latest

	`, projectName, projectName)

	e := os.WriteFile("Makefile", []byte(baseMakefile), 0755)
	if e != nil {
		fmt.Printf("Unable to create Dockerfile: %v", e)
	}
}
