.PHONY: test
cstar: main.go commands.go
	go fmt && go build
test:
	go fmt && go test
