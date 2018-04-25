PLUGIN := terraform-provider-lambda 


$(PLUGIN): vet fmt test
	go build -o $@ .

.PHONY: vet
vet:
	go $@ ./...

.PHONY: fmt
fmt:
	go $@ ./...

.PHONY: test
test:
	go $@ -v ./...
