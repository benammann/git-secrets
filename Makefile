watch:
	CompileDaemon --build "go build -o bin/git-secrets ."

tests:
	go test ./...