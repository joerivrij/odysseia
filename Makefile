git_short_hash=`git rev-parse --short HEAD`
project_name='odysseia'
image := $(shell docker images -q ${project_name}:${git_short_hash})

create-image:
ifeq ("${image}","")
	echo "creating base image"
	echo "docker build ${project_name}:$(git_short_hash)"
	docker build -t $(project_name):$(git_short_hash) . --no-cache
else
	echo "${image}"
	echo "base image already present on this machine"
endif


