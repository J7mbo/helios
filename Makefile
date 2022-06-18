build: build-highlighter
	go build -o run.bin main.go

build-highlighter:
	go build -o highlighter/highlighter.bin highlighter/highlighter.go