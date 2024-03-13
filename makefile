install:
	go mod tidy
	bun install

dev:
	chmod +x dev.sh && ./dev.sh