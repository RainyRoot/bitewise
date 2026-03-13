.PHONY: dev dev-backend dev-mobile build migrate clean

dev:
	@echo "Starting BiteWise development servers..."
	@trap 'kill 0' INT TERM; \
	$(MAKE) dev-backend & \
	$(MAKE) dev-mobile & \
	wait

dev-backend:
	cd backend && go run ./cmd/server/

dev-mobile:
	cd mobile && export PATH="$$HOME/.local/bin:$$PATH" && pnpm start

build:
	cd backend && go build -o bin/server ./cmd/server/

migrate:
	cd backend && go run ./cmd/server/ -migrate

clean:
	rm -rf backend/bin/
