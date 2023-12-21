default: gen lint test build

# Run acceptance tests
.PHONY: testacc
testacc: gen lint
	# INFO: acceptance tests use your api_token.
	# Consider setting METAL_STACK_CLOUD_API_TOKEN.
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Run unit tests
.PHONY: test
test:
	go test ./... -v $(TESTARGS) -timeout 120m

# Generate docs
.PHONY: gen
gen:
	go mod tidy
	go generate

# Build
.PHONY: build
build:
	go build .

# Lint
.PHONY: lint
lint:
	golangci-lint run

# Check Licenses
.PHONY: check-licenses
check-licenses:
	# Requires go install github.com/google/go-licenses@latest
	go-licenses check --ignore github.com/metal-stack-cloud --include_tests .
