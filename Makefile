.PHONY: gen/test
gen/test:
	cd test && go generate

.PHONY: test
test:
	cd test && deno test --allow-read .

.PHONY: fmt
fmt:
	deno fmt . ./test