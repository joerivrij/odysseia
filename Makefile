docker-build:
	cd sokrates && make docker-build && cd ..
	cd pheidias && make docker-build