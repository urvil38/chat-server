clean:
	rm -rf ./bin

linux-build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/chat-client -ldflags="-w" ./client
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/chat-server -ldflags="-w" ./server

macos-build:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/chat-client -ldflags="-w" ./client
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/chat-server -ldflags="-w" ./server

windows-build:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/chat-client.exe -ldflags="-w" ./client
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/chat-server.exe -ldflags="-w" ./server
