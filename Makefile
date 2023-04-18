build/bin/golangci-lint: Makefile
	mkdir -p build/bin
	# Linux tar requires --wildcards for universal pattern. Darwin doesn't have such option.
	curl -L https://github.com/golangci/golangci-lint/releases/download/v1.52.2/golangci-lint-1.52.2-`go \
		env GOHOSTOS`-`go env GOHOSTARCH`.tar.gz | tar -C build/bin --strip-components 1 \
		`[[ \`uname -s\` == Linux ]] && echo --wildcards` -zx \*/golangci-lint
	touch -c build/bin/golangci-lint

lint: build/bin/golangci-lint
	build/bin/golangci-lint run ./...
.PHONY: lint

test:
	go test -count=1 -v ./...
.PHONY: test
