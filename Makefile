SRC = ./src/...
COVER = cover/cover

run:
	go run .

build-linux:
	GOOS=linux GOARCH=amd64 go build  -ldflags "-s -w" -o bin/parse main.go

build-mac:
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o bin/parse main.go

test:
	go test -p 1 $(SRC)

test-cover:
	go test -p 1 $(SRC) -v -coverprofile $(COVER).out && go tool cover -html $(COVER).out -o $(COVER).html && open $(COVER).html