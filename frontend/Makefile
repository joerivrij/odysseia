git_short_hash=`git rev-parse --short HEAD`

docker-build:
	docker build -t lexiko_spa:$(git_short_hash) .

docker-run: docker-build
	 docker run -d -p 5000:5000 --name lexiko_spa lexiko_spa:$(git_short_hash)

docker-remove:
	docker rm -f lexiko_spa