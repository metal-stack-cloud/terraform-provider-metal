default: gen lint test build

# Run acceptance tests
.PHONY: testacc
testacc: gen lint
	# ATTENTION: acceptance tests run against metalstack.cloud by default!
	# Consider setting METAL_STACK_CLOUD_API_URL, METAL_STACK_CLOUD_API_TOKEN, METAL_STACK_CLOUD_ORGANIZATION, METAL_STACK_CLOUD_PROJECT.
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
