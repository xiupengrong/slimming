.PHONY: build test clean

build:
	go build -o slimming.exe .

test:
	go test ./...

clean:
	rm -f slimming.exe
