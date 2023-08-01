SRC = ./src/...
COVER = cover/cover
LINUX32 = GOOS=linux GOARCH=386
LINUX64 = GOOS=linux GOARCH=amd64
MAC32 = GOOS=darwin GOARCH=386
MAC64 = GOOS=darwin GOARCH=amd64

run:
	go run .

build-linux32:
	$(LINUX32) go build -o bin/parse src/parse/main.go && $(LINUX32) go build -o bin/beautify src/beautify/main.go

build-linux:
	$(LINUX64) go build -o bin/parse src/parse/main.go && $(LINUX64) go build -o bin/beautify src/beautify/main.go

build-mac32:
	$(MAC32) go build -o bin/parse src/parse/main.go && $(MAC32) go build -o bin/beautify src/beautify/main.go

build-mac:
	$(MAC64) go build -o bin/parse src/parse/main.go && $(MAC64) go build -o bin/beautify src/beautify/main.go

test:
	go test -p 1 $(SRC)

test-cover:
	go test -p 1 $(SRC) -v -coverprofile $(COVER).out && go tool cover -html $(COVER).out -o $(COVER).html && open $(COVER).html