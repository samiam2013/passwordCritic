build: 
	go build ./passchk/main.go

test: 
	go test ./... # --race 

coverage:
	bash ./codecov.sh

run: 
	go run ./passchk/main.go

.PHONY: test build run