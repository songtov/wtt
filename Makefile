VERSION ?= v0.6.1

build:
	go build -ldflags "-X github.com/songtov/wtt/cmd.Version=$(VERSION)" -o wtt-bin .

test-local: build
	@echo ""
	@echo "Local binary built. To test in this shell:"
	@echo "  export PATH=\$$(pwd):\$$PATH"
	@echo "  eval \"\$$(wtt-bin --init zsh)\""
	@echo "  wtt version  # should show $(VERSION)"
	@echo ""

push-tag:
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "v0.3.0-dev" ]; then \
		echo "Error: set a real VERSION, e.g. make push-tag VERSION=v0.3.0"; \
		exit 1; \
	fi
	git tag $(VERSION)
	git push origin $(VERSION)
	@echo "Tag $(VERSION) pushed — GitHub Actions release workflow triggered."

clean:
	rm -f wtt-bin
