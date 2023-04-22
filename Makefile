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

RUN_COUNT ?= 6
CPU ?= 1

bench:
ifdef out_file
	go test -bench=. -benchmem -cpu ${CPU} -count $(RUN_COUNT) -run=^# ./tests/benchmark/*_test.go > $(out_file)
else
	go test -bench=. -benchmem -cpu ${CPU} -count $(RUN_COUNT) -run=^# ./tests/benchmark/*_test.go
endif
.PHONY: bench

old_bench.out:
	git stash; out_file=old_bench.out make bench; git stash pop; echo "OK"

bench-cmp: old_bench.out
	out_file=new_bench.out make bench
	benchstat old_bench.out new_bench.out
.PHONY: bench-cmp

bench-base:
	out_file=old_bench.out make bench