
test: verify
	go get  ./streams/...
	go test ./streams/...

verify:
	golint  ./streams/...
	go vet  ./streams/...

.PHONY: test verify
