all: lint test
.PHONY: all

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
	go test -count=1 --race -v ./... -coverprofile=coverage.out
.PHONY: test

test-cpu-profile:
	go test ./...
.PHONY: test-cpu-profile

bench:
	go test -bench=./...
.PHONY: bench

old_bench.out:
	git stash
	go test -bench=. -count 5 ./query/builder_test.go > old_bench.out
	git stash pop

bench-cmp: old_bench.out
	go test -bench=. -count 5 ./query/builder_test.go > new_bench.out
	benchcmp old_bench.out new_bench.out
.PHONY: bench-cmp