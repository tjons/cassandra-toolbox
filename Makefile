default: test vet fmt

bench:
	go test -bench=. -benchmem -count 10 ./...

.PHONY: test
test:
	go test -race -cover ./...

vet:
	go vet ./...

fmt:
	go fmt ./...