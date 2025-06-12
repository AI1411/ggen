# PostgreSQL connection settings
DB_HOST ?= postgres
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= gen
DB_SSLMODE ?= disable
# PostgreSQL connection settings for test database
DB_TEST_HOST ?= postgres-test
DB_TEST_PORT ?= 5432
DB_TEST_USER ?= postgres
DB_TEST_PASSWORD ?= postgres
DB_TEST_NAME ?= gen_test
DB_TEST_SSLMODE ?= disable

# Construct database URL
DATABASE_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
DATABASE_TEST_URL = postgres://$(DB_TEST_USER):$(DB_TEST_PASSWORD)@$(DB_TEST_HOST):$(DB_TEST_PORT)/$(DB_TEST_NAME)?sslmode=$(DB_TEST_SSLMODE)

.PHONY: migrate migrate-up migrate-down migrate-version migrate-create

# Run migration up
migrate:
	docker compose exec migrations migrate -source file://./ -database '$(DATABASE_URL)' up
	docker compose exec migrations migrate -source file://./ -database '$(DATABASE_TEST_URL)' up

# Alias for migrate
migrate-up: migrate

# Run migration down
migrate-down:
	docker compose exec migration migrate -source file://./ -database '$(DATABASE_URL)' down 1
	docker compose exec migration migrate -source file://./ -database '$(DATABASE_TEST_URL)' down 1

# Show current migration version
migrate-version:
	docker compose exec migration migrate -source file://./ -database '$(DATABASE_URL)' version
	docker compose exec migration migrate -source file://./ -database '$(DATABASE_TEST_URL)' version

# Create new migration file
migrate-create:
	@read -p "Enter migration name: " name; \
	docker compose exec migration migrate create -ext sql -dir ./ -seq $$name

.PHONY: logs
logs:
	docker logs api -f --tail 100

.PHONY: generate-models
generate-models:
	@docker compose exec api go run ./cmd/gormgen/generate_all/main.go
	@docker compose exec api go run ./cmd/gormgen/generate_associations/main.go

.PHONY: exec-schema
exec-schema: ## sqlファイルをコンテナに流す
	cat ./backend/migrations/*.up.sql > ./backend/migrations/schema.sql
	docker cp backend/migrations/schema.sql postgres:/schema.sql
	docker cp backend/migrations/schema.sql postgres-test:/schema.sql
	docker exec -it postgres psql -U postgres -d gen -f /schema.sql
	docker exec -it postgres-test psql -U postgres -d gen_test -f /schema.sql
	rm ./backend/migrations/schema.sql

.PHONY: swag
swag: ## swagger更新
	@docker compose exec api swag init -g ./cmd/api/main.go --output ./docs/api
	@cd frontend && pnpm generate
fmt: ## コードを自動整形（ツールチェイン使用）
	@cd backend && go run mvdan.cc/gofumpt@latest -l -w .
	@cd backend && go run golang.org/x/tools/cmd/goimports@latest -l -w -local "g_gen" .
	@cd frontend && pnpm format

.PHONY: lint lint-fix test test-coverage vet sec staticcheck tools
## 開発ツール関連

tools: ## 開発ツールをインストール（ツールチェイン使用）
	@cd backend && go mod download
	@cd backend && go mod tidy

lint: ## 静的解析実行（ツールチェイン使用）
	@cd backend && go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run ./...

lint-fix: ## 静的解析で修正可能な問題を自動修正（ツールチェイン使用）
	@cd backend && go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --fix ./...

vet: ## go vetを実行
	@cd backend && go vet ./...

staticcheck: ## staticcheckを実行（ツールチェイン使用）
	@cd backend && go run honnef.co/go/tools/cmd/staticcheck@latest ./...

sec: ## セキュリティチェック実行（ツールチェイン使用）
	@cd backend && go run golang.org/x/vuln/cmd/govulncheck@latest ./...

test: ## テスト実行
	docker compose exec api gotestsum --format testname ./...

test-coverage: ## テストカバレッジ計測
	@cd backend && go test -coverprofile=coverage.out ./...
	@cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "カバレッジレポートが coverage.html に生成されました"

quality: lint vet staticcheck sec ## コード品質チェック（全ツール）

ci: fmt quality test ## CI環境で実行するチェック

mockgen: ## モックを生成
	docker compose exec api go generate ./...

seeder: ## シーダーを実行
	docker compose exec api go run ./cmd/seed/municipality/main.go
.PhONY: tidy
tidy: ## 依存関係の整理
	docker compose exec api go mod tidy