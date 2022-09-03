.PHONY: clean build packing

clean: 
	rm -rf ./main
	
build:
	@GOOS=linux GOARCH=amd64
	@echo ">> Building GRPC..."
	@go build -o service-area-grpc ./cmd/grpc
	@echo ">> Finished"

build-ci:
	@GOOS=linux GOARCH=amd64
	@echo ">> Building GRPC..."
	@CGO_ENABLED=0  go build -o ./bin/service-area-grpc ./cmd/grpc
	@echo ">> Finished"

packing:
	zip service-area-grpc.zip ./cmd/grpc

run:
	./service-area-grpc

test:
	ginkgo -r

test-local:
	go test ./... -coverprofile cover.out.temp && cat cover.out.temp | grep -v 'mock_*\|\.pb\.' > cover.out && go tool cover -func cover.out
	
mock-generate:
	@mockgen -source=repository/interface.go -destination=repository/mock_repository.go -package=repository
