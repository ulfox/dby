.PHONY: clean
clean:
	rm -rf .test

.PHONY: clean
lint:
	@if ! command -v golint; then \
		go get -u golang.org/x/lint/golint; \
	fi
	golint -set_exit_status ./...


.PHONY: test
test:
	go test -v .
