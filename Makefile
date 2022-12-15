EXE=orquestrator

all: build

build:
	go mod init puppy
	go build -o bin/$(EXE) -modfile source/go.mod source/main.go source/messages.go source/sessionControl.go
	rm -f go.mod

clean:
	rm -f bin/$(EXE)
	