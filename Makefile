EXE=orquestrator

all: build

build:
	go mod init tidy
	go build -o bin/$(EXE) source/main.go

clean:
	rm -f bin/$(EXE)
	