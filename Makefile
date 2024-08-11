BINARY_NAME=syllabuster
 
all: build
 
build:
	mkdir -p ./bin
	go build -o ./bin/${BINARY_NAME} main.go
test:
	go test -v main.go
run:
	mkdir -p ./bin
	go build -o ./bin/${BINARY_NAME} main.go
	./bin/${BINARY_NAME}
clean:
	go clean
	rm ./bin/${BINARY_NAME}
