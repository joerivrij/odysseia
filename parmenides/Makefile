git_short_hash=`git rev-parse --short HEAD`
project_name='parmenides'

build:
	GO111MODULE=on GOOS=linux CGO_ENABLED=0 go build -o $(project_name) cmd/main.go
	mkdir ./bin
	mv $(project_name) ./bin/

docker-build:
	docker build -t $(project_name):$(git_short_hash) .
