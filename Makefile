GOPATH=$(shell go env GOPATH)

generate:
	GOBIN=$(GOPATH)/bin go install github.com/ogen-go/ogen/cmd/ogen@v0.76.0
	$(GOPATH)/bin/ogen --target admin/generated -package generated --no-client --convenient-errors=on --clean admin/openapi.yaml

