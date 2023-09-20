BINARY_NAME=onionsweep

GOBUILD=go build
GOCLEAN=go clean

MAIN_PATH=cmd/main.go

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)