.PHONY: default lint test yaegi_test vendor clean

export GO111MODULE=on

default: lint test yaegi_test

lint:
	golangci-lint run

test:
	go test -v -cover ./...

# Runs the tests through the Yaegi interpreter — the same way Traefik loads the
# plugin at runtime. Passing `go test` does not guarantee Yaegi compatibility.
yaegi_test:
	yaegi test -v .

vendor:
	go mod vendor

clean:
	rm -rf ./vendor
