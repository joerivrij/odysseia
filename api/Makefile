git_short_hash=`git rev-parse --short HEAD`

docker-build:
	docker build -t lexico_api:$(git_short_hash) .

docker-run: docker-build
	 docker run -d -p 5000:5000 --name lexico_api lexico_api:$(git_short_hash)

docker-remove:
	docker rm -f lexico_api