git_short_hash=`git rev-parse --short HEAD`
current_dir=`pwd`
project_name := $(shell basename ${current_dir})
base_name='odysseia'
harbor_address='core.harbor.domain:30003'

setup:
	if [ -e ./dist ]; then rm -rf ./dist; fi

dist:
	mkdir -p ./dist
	mkdir -p ./dist/bin

build: test setup dist
	go fmt ./...
	cd ..; GO111MODULE=on GOOS=linux CGO_ENABLED=0 go build $(project_name)/main.go;mv main ./$(project_name)/dist/bin/$(project_name)

test:
	go test ./... -cover

test-with-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

copy-dockerfile:
	cp ../api.Dockerfile .

copy-local-dockerfile:
	cp ../local.Dockerfile .

create-base-image:
	cd ..; make create-image

create-image: test create-base-image copy-dockerfile
	echo "docker build ${project_name}:$(git_short_hash)"
	cd ..; docker build --build-arg project_name=$(project_name) --build-arg git_hash=$(git_short_hash) -f $(project_name)/api.Dockerfile -t $(project_name):$(git_short_hash) . --no-cache

create-local-image: copy-local-dockerfile build
	echo "docker build ${project_name}"
	cp -r dist/bin/$(project_name) .
	docker build --build-arg project_name=$(project_name) -f local.Dockerfile -t $(project_name):$(git_short_hash) . --no-cache
	rm $(project_name)

create-harbor: create-local-image
	docker tag $(project_name):$(git_short_hash) $(harbor_address)/$(base_name)/$(project_name):$(git_short_hash)
	docker push $(harbor_address)/$(base_name)/$(project_name):$(git_short_hash)