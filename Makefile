build:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -o bin/listen

run:
	./bin/listen

clean:
	rm -rf ./bin ./vendor go.sum
