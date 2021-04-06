git_short_hash=`git rev-parse --short HEAD`

build:
	GO111MODULE=on GOOS=linux CGO_ENABLED=0 go build -o sokrates cmd/main.go
	mkdir ./bin
	mv sokrates ./bin/

docker-build:
	docker build -t sokrates:$(git_short_hash) .
