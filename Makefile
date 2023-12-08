test: verify
	go test ./internal/...

verify:
	staticcheck ./internal/...
	go vet ./internal/...

.PHONY: test verify
