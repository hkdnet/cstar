.PHONY: test
.PHONY: run
cstar: main.go commands.go command/list.go
	go fmt && go vet && go build
test:
	go fmt && go vet && go test ./command
run:
	go fmt && go vet && go build && ./cstar
build: main.go commands.go command/list.go
	mkdir -p release/win/
	GOOS=windows GOARCH=amd64 go build && mv -f cstar.exe release/win/cstar-x64.exe
	GOOS=windows GOARCH=386   go build && mv -f cstar.exe release/win/cstar-x86.exe
	mkdir -p release/osx/
	GOOS=darwin GOARCH=amd64  go build && mv -f cstar release/osx/cstar
