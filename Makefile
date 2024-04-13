export DB_USER=book_db
export DB_NAME=book_db
export DB_PASSWORD=book_db
export DB_HOST=book-db
export DB_PORT=3306
export DB_SSL_MODE=disable

PWD?=$(shell pwd)

deps:
	go install github.com/golang/mock/mockgen@v1.6.0 && \
	go mod download

build:
	go build -o /book-service

test:
	go test -timeout=120s -cover -race ./...

clean_mocks:
	rm -rf mocks/

mocks: clean_mocks
	go generate ./...

docker_build:
	docker build --ssh default . \
		-t book_service

docker_test: docker_build
	docker run \
		--rm \
		-v $(PWD):$(PWD) -w $(PWD) -v /var/run/docker.sock:/var/run/docker.sock \
		book_service \
		/usr/bin/make test
	docker --log-level ERROR compose -p book_test stop
