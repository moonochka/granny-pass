GO_CMD=CGO_ENABLED=0 GO111MODULE=on time go

test:
	$(GO_CMD) test -tags graphTest ./...
	$(GO_CMD) test -tags processorTest ./...

run:
	$(GO_CMD) run cmd/granny-pass-dev/main.go -k

run-full:
	$(GO_CMD) run cmd/granny-pass-dev/main.go -min 20 -max 24 -cnt 4 -k -file 40000.txt

run-full-task:
	$(GO_CMD) run cmd/granny-pass-dev/main.go -min 20 -max 24 -cnt 4 -file 40000.txt