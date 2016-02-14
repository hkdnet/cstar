.PHONY: test
.PHONY: run
cstar: main.go commands.go command/list.go
	go fmt && go vet && go build
test:
	go fmt && go vet && go test
run:
	go fmt && go vet && go build && ./cstar
