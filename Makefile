.PHONY: cmd
start:
	quayd -github-token=$GITHUB_TOKEN -port=$PORT -registry-auth=$REGISTRY_AUTH
cmd:
	godep go build -o build/quayd ./cmd/quayd
