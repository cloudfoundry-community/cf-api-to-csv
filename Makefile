.PHONY: darwin linux
darwin:
	env GOOS=darwin GOARCH=amd64 go build
linux:
	env GOOS=linux GOARCH=amd64 go build