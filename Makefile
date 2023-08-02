SRC = ./src/...
COVER = cover/cover
LINUX32 = GOOS=linux GOARCH=386
LINUX64 = GOOS=linux GOARCH=amd64
MAC64 = GOOS=darwin GOARCH=arm64

run:
	go run .

build-linux:
	$(LINUX64) go build -o bin/parse main.go

build-mac:
	$(MAC64) go build -o bin/parse main.go

test:
	go test -p 1 $(SRC)

test-cover:
	go test -p 1 $(SRC) -v -coverprofile $(COVER).out && go tool cover -html $(COVER).out -o $(COVER).html && open $(COVER).html