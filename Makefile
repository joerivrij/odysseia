docker-build:
	cd api && make docker-build && cd ..
	cd frontend && make docker-build