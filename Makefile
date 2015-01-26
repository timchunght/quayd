.PHONY: cmd

cmd:
	godep go build -o build/quayd ./cmd/quayd
