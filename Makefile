VERSION ?= v0.3.0-dev

build:
	go build -ldflags "-X github.com/songtov/wtt/cmd.Version=$(VERSION)" -o wtt-bin .

test-local: build
	@echo ""
	@echo "Local binary built. To test in this shell:"
	@echo "  export PATH=$(PWD):$$PATH"
	@echo "  eval \"\$$(wtt-bin --init zsh)\""
	@echo "  wtt version  # should show $(VERSION)"
	@echo ""

clean:
	rm -f wtt-bin
