git_short_hash=`git rev-parse --short HEAD`
project_name='hippokrates'

docker-build:
	echo "docker build ${project_name}"
	docker build -t $(project_name):$(git_short_hash) . --no-cache
