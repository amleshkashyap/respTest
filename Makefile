build:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -o listen

run:
	./listen "mode-2"

runCustom:
	./listen "mode-2" "Custom"

runSample:
	./listen "mode-2" "samples"

runSampleCustom:
	./listen "mode-2" "samples" "Custom"

runMode1:
	./listen "mode-1"

clean:
	rm -rf ./bin ./vendor go.sum
