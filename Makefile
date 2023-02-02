EXE=orquestrator

all: build

build:
	go mod tidy
	go build -o bin/$(EXE) main.go

clean:
	rm -f bin/$(EXE)
	