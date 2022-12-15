EXE=worker

all: build

build:
	go mod init puppy
	go build -o bin/$(EXE) -modfile source/go.mod source/main.go source/messages.go source/experimentControl.go
	rm -f go.mod

run: build
	./bin/worker --timeout 50

clean:
	rm -f bin/$(EXE)
	