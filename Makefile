.PHONY: test

test:
	go test -v -coverprofile=coverage.out --race ./...
	go tool cover -func=coverage.out

docker-test:
	docker build -t docker-test .
	docker run docker-test make
