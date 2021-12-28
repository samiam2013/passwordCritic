rebuild: 
	wget "https://raw.githubusercontent.com/danielmiessler/SecLists/master/Passwords/Common-Credentials/10-million-password-list-top-1000000.txt" \
		-o ./cache/top_10MM_passwords.txt;

build: 
	go build ./passchk/main.go

test: 
	go test ./... # --race 

coverage:
	bash ./codecov.sh

run: 
	go run ./passchk/main.go

.PHONY : rebuild test build run