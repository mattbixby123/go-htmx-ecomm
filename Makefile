.PHONY: run build clean dev

# Run the app
run:
	go run .

# Build the binary
build:
	go build -o ecommerce

# Clean build artifacts
clean:
	rm -f ecommerce

# Development mode with auto-reload (if you install air)
dev:
	air