.PHONY: test
.PHONY: run
.PHONY: release
cstar: main.go commands.go command/list.go
	go fmt && go vet && go build
test:
	go fmt && go vet && go test ./command
run:
	go fmt && go vet && go build && ./cstar
release:
	mkdir -p release/win/x64
	GOOS=windows GOARCH=amd64 go build && mv -f cstar.exe release/win/x64/cstar.exe
	mkdir -p release/win/x86
	GOOS=windows GOARCH=386   go build && mv -f cstar.exe release/win/x86/cstar.exe
	mkdir -p release/osx/
	GOOS=darwin GOARCH=amd64  go build && mv -f cstar release/osx/cstar
