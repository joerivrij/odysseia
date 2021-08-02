git_short_hash=`git rev-parse --short HEAD`
project_name=xerxes

create-image:
	echo "docker build ${project_name}:$(git_short_hash)"
	docker build -t $(project_name):$(git_short_hash) .
