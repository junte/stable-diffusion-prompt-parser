SRC = ./src/...
COVER = cover/cover

run:
	go run .

build:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/linux-amd64/parse main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o bin/mac-arm64/parse main.go

test:
	go test -p 1 $(SRC)

test-cover:
	go test -p 1 $(SRC) -v -coverprofile $(COVER).out && go tool cover -html $(COVER).out -o $(COVER).html && open $(COVER).html