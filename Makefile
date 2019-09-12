vet:
	@go vet ./...

test: vet
	@go test ./... -coverprofile cover.out
