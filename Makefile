rebuild: 
	cd ./cache/; \
	wget "https://raw.githubusercontent.com/danielmiessler/SecLists/master/Passwords/Common-Credentials/10-million-password-list-top-1000000.txt"; \
	wget "https://raw.githubusercontent.com/danielmiessler/SecLists/master/Passwords/Common-Credentials/10-million-password-list-top-100000.txt"; \
	wget "https://raw.githubusercontent.com/danielmiessler/SecLists/master/Passwords/Common-Credentials/10-million-password-list-top-10000.txt"; \
	wget "https://raw.githubusercontent.com/danielmiessler/SecLists/master/Passwords/Common-Credentials/10-million-password-list-top-1000.txt";

build: 
	go build ./passchk/main.go

test: 
	go test ./... # --race 

coverage:
	bash ./codecov.sh

run: 
	go run ./passchk/main.go

.PHONY : rebuild test build run