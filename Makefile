SHELL:=/bin/bash

.PHONY: clean
clean:
	rm -rf tests/.test;

.PHONY: lint
lint:
	@if ! command -v golint; then \
		go get -u golang.org/x/lint/golint; \
	fi
	golint -set_exit_status ./...


.PHONY: test
test: clean
	@set -e; \
	hasErr=0; \
	for i in {0..10}; do \
		go test -v ./tests/; \
		if [[ "$${?}" != 0 ]]; then \
			hasErr=1; break; \
		fi; \
	done; \
	rmdir ./tests/.test
