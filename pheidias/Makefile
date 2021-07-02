git_short_hash=`git rev-parse --short HEAD`
project_name=pheidias

lint:
	npm run lint

create-image: lint
	echo "docker build ${project_name}:$(git_short_hash)"
	docker build -t $(project_name):$(git_short_hash) .
