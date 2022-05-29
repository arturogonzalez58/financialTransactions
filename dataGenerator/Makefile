build:
	go build -o bin/dataGenerator cmd/main.go

run:
	go run cmd/main.go

build-docker-image:
	docker image build -t data-generator .

run-docker:
	docker run --env-file env data-generator