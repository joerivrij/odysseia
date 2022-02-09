git_short_hash=`git rev-parse --short HEAD`
project_name=xerxes
base_name='odysseia'
harbor_address='core.harbor.domain:30003'

create-image:
	echo "docker build $(project_name):$(git_short_hash)"
	docker build -t $(project_name):$(git_short_hash) .

create-harbor: create-image
	docker tag $(project_name):$(git_short_hash) $(harbor_address)/$(base_name)/$(project_name):$(git_short_hash)
	docker push $(harbor_address)/$(base_name)/$(project_name):$(git_short_hash)