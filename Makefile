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
	@go test -v ./tests/; \
	rmdir ./tests/.test

.PHONY: testloop
testloop: clean
	@set -e; \
	hasErr=0; \
	for i in $$(seq 0 $${test_loop}); do \
		go test -v ./tests/; \
		if [[ "$${?}" != 0 ]]; then \
			hasErr=1; break; \
		fi; \
	done; \
	rmdir ./tests/.test
