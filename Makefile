BINARY_NAME=gender-engine

# remove any binaries that are built
clean:
	rm -f ./bin/$(BINARY_NAME)*

build-debug: clean
	CGO_ENABLED=0 go build -gcflags=all="-N -l" -o bin/$(BINARY_NAME)-debug main.go
