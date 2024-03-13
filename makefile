install:
	go mod tidy
	bun install
	bun run dev:css

dev:
	air -c ./.air.toml