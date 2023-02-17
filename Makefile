EXE=worker

all: build

build:
	go mod tidy
	go build -o bin/$(EXE) main.go

run: build
	./bin/worker --timeout 50

clean:
	rm -f bin/$(EXE)
	