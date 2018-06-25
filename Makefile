.PHONY: darwin linux
darwin:
	env GOOS=darwin GOARCH=amd64 go build -o cf-metrics-darwin
linux:
	env GOOS=linux GOARCH=amd64 go build -o cf-metrics-linux