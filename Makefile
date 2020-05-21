
test: verify
	go test ./streams/...

verify:
	golint ./streams/...
	go vet ./streams/...

.PHONY: test verify
