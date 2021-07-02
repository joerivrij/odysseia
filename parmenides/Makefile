git_short_hash=`git rev-parse --short HEAD`
project_name='parmenides'

setup:
	if [ -e ./dist ]; then rm -rf ./dist; fi

dist:
	mkdir -p ./dist
	mkdir -p ./dist/linux/bin

build: setup dist
	go fmt ./...
	cd ..; GO111MODULE=on GOOS=linux CGO_ENABLED=0 go build $(project_name)/main.go;mv main ./$(project_name)/dist/linux/bin/$(project_name)

create-image: build
	echo "docker build ${project_name}:$(git_short_hash)"
	cp -r dist/linux/bin/$(project_name) .
	docker build -t $(project_name):$(git_short_hash) . --no-cache
	rm $(project_name)