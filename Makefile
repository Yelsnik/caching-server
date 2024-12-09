.PHONY: build

install:
	go build
	go install

run:
	go build
	go install
	caching-server --origin https://trackinginventory.onrender.com --port 8080