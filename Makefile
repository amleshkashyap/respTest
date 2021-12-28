build:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -o listen

run:
	./listen "mode-2"

runSample:
	./listen "mode-2" "value"

runMode1:
	./listen "mode-1"

clean:
	rm -rf ./bin ./vendor go.sum
