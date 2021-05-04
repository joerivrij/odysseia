git_short_hash=`git rev-parse --short HEAD`
project_name='pheidias'

docker-build:
	docker build -t $(project_name):$(git_short_hash) .