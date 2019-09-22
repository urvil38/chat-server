linux-build:
	GOOS=linux GOARCH=amd64 go build -o chat-server -ldflags="-w" ./server
	GOOS=linux GOARCH=amd64 go build -o chat-client -ldflags="-w" ./client

macos-build:
	GOOS=darwin GOARCH=amd64 go build -o chat-server -ldflags="-w" ./server
	GOOS=darwin GOARCH=amd64 go build -o chat-client -ldflags="-w" ./client

windows-build:
	GOOS=windows GOARCH=amd64 go build -o chat-server.exe -ldflags="-w" ./server
	GOOS=windows GOARCH=amd64 go build -o chat-client.exe -ldflags="-w" ./client
