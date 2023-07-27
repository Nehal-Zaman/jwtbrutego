build: 
	@go build -o bin/jwtbrutego  -ldflags="-s -w" 

run: build
	@./bin/jwtbrutego