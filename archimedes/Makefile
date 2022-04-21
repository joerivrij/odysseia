project_name='archimedes'

build:
	go fmt ./...
	GO111MODULE=on GOOS=linux CGO_ENABLED=0 go build cmd/archimedes.go;
	sudo mv archimedes /usr/local/bin/$(project_name)

build-mac:
	go fmt ./...
	GO111MODULE=on GOOS=darwin CGO_ENABLED=0 go build cmd/archimedes.go;
	mv archimedes /usr/local/bin/$(project_name)