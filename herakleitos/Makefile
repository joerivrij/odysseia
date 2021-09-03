git_short_hash=`git rev-parse --short HEAD`
current_dir=`pwd`
project_name := $(shell basename ${current_dir})
base_name='odysseia'

setup:
	if [ -e ./dist ]; then rm -rf ./dist; fi

dist:
	mkdir -p ./dist
	mkdir -p ./dist/linux/bin

build: test setup dist
	go fmt ./...
	cd ..; GO111MODULE=on GOOS=linux CGO_ENABLED=0 go build $(project_name)/main.go;mv main ./$(project_name)/dist/linux/bin/$(project_name)

test:
	go test ./... -cover

test-with-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

create-base-image:
	cd ..; make create-image

create-image: test create-base-image
	echo "docker build ${project_name}:$(git_short_hash)"
	cd ..; docker build --build-arg project_name=$(project_name) --build-arg git_hash=$(git_short_hash) -f $(project_name)/Dockerfile -t $(project_name):$(git_short_hash) . --no-cache
