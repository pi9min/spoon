TEST_LIST = $(shell go list ./...)

test: dep lint
	go test -v $(TEST_LIST) -count=1

lint: dep
	go vet ./...

dep:
	go mod download
