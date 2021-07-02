git_short_hash=`git rev-parse --short HEAD`
project_name='herodotos'

setup:
	if [ -e ./dist ]; then rm -rf ./dist; fi

dist:
	mkdir -p ./dist
	mkdir -p ./dist/linux/bin
#	mkdir -p ./dist/arm64/bin

test:
	go test ./... -cover

build: test setup dist
	go fmt ./...
	cd ..; GO111MODULE=on GOOS=linux CGO_ENABLED=0 go build $(project_name)/cmd/main.go;mv main ./$(project_name)/dist/linux/bin/$(project_name)
	#cd ..; GO111MODULE=on GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(project_name)/cmd/main.go;mv main ./$(project_name)/dist/arm64/bin/$(project_name)

create-image: build
	echo "docker build ${project_name}:$(git_short_hash)"
	cp -r dist/linux/bin/$(project_name) .
	docker build -t $(project_name):$(git_short_hash) . --no-cache
	rm $(project_name)
