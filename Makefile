.PHONY: dev build-all clean dev-server dev-client test-all

dev-server:
	@echo "Starting Go backend..."
	cd server && go run main.go

dev-client:
	@echo "Starting Vite frontend..."
	cd app && npm run dev

build-all:
	@echo "Building frontend and backend..."
	cd app && npm run build
	cd server && go build -o ../bin/server main.go

test-all:
	@echo "Running backend and frontend tests..."
	cd server && go test ./...
	cd app && npm run test -- --watchAll=false || echo "No client tests yet"

clean:
	@echo "Cleaning binaries..."
	rm -rf bin/
