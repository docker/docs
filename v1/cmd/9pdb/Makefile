all: 9pdb.exe 9pdb

.PHONY: 9pdb.exe
9pdb.exe:
	GOOS=windows GOARCH=amd64 go build -o 9pdb.exe

.PHONY: 9pdb
9pdb:
	go build
