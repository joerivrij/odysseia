git_short_hash=`git rev-parse --short HEAD`
project_name='hippokrates'

copy-go-mod:
	cp ../go.mod .
	cp ../go.sum .

create-image:
	echo "docker build ${project_name}:$(git_short_hash)"
	docker build -t $(project_name):$(git_short_hash) . --no-cache
