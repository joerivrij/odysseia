git_short_hash=`git rev-parse --short HEAD`
project_name='ploutarchos'

docker-build:
	echo "docker build ${project_name}"
	docker build -t $(project_name):$(git_short_hash) . --no-cache


build-local:
	docker rm -f $(project_name)
	docker build -t $(project_name):$(git_short_hash) . --no-cache
	docker run -d -p 8080:8080 --name $(project_name) $(project_name):$(git_short_hash)
