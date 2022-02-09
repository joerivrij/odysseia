git_short_hash=`git rev-parse --short HEAD`
project_name='ploutarchos'
base_name='odysseia'
harbor_address='core.harbor.domain:30003'

create-image:
	echo "docker build $(project_name):$(git_short_hash)"
	docker build -t $(project_name):$(git_short_hash) . --no-cache

build-local:
	docker rm -f $(project_name)
	docker build -t $(project_name):$(git_short_hash) . --no-cache
	docker run -d -p 8080:8080 --name $(project_name) $(project_name):$(git_short_hash)

create-harbor: create-image
	docker tag $(project_name):$(git_short_hash) $(harbor_address)/$(base_name)/$(project_name):$(git_short_hash)
	docker push $(harbor_address)/$(base_name)/$(project_name):$(git_short_hash)