
test:
	go test -race -coverprofile=coverage.out ./...

cover: test
	go tool cover -html=coverage.out
