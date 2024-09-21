build:
	go build -o ./bin/blockchainsystem

run: build 
		./bin/blockchainsystem


test: 
	go test  ./...