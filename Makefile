client:
	air --build.cmd "go build -o ./bin/client ./cmd/client/main.go" --build.bin "./bin/client -folder ../tmp"

server:
	air --build.cmd "go build -o ./bin/server ./cmd/server/main.go" --build.bin ./bin/server