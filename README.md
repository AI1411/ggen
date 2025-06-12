## 環境構築

### 前提条件

- Docker
- Docker Compose
- Node.js (v16以降) - フロントエンド開発時のみ
- pnpm - フロントエンド開発時のみ

### 開発環境の起動

1. プロジェクトをクローン
```bash
git clone <repository-url>
cd g_gen
```

2. 環境変数ファイルを設定
```bash
cp .env.example .env  # 必要に応じて値を編集
```

3. Dockerコンテナの起動
```bash
docker compose up -d
```

4. データベースマイグレーション実行
```bash
make migrate
```

### バックエンド開発

#### 技術スタック
- Go
- PostgreSQL
- Docker/Docker Compose

#### 環境構築手順

1. バックエンドディレクトリに移動
```bash
cd backend
```

2. 依存関係のダウンロード
```bash
go mod download
```

3. 開発サーバーの起動
Docker Composeでの起動により、ホットリロード機能付きでサーバーが開始されます
- API: http://localhost:8080
- Swagger UI: http://localhost:8084

#### 開発用コマンド

```bash
# コード整形
make fmt

# 静的解析
make lint

# テスト実行
make test

# モック生成
make mockgen

# Swagger更新
make swag
```

### フロントエンド開発

#### 技術スタック
- React
- TypeScript
- Vite
- Tailwind CSS
- TanStack Query
- TanStack Router

#### 環境構築手順

1. フロントエンドディレクトリに移動
```bash
cd frontend
```

2. 依存関係のインストール
```bash
pnpm install
```

3. 開発サーバーの起動
```bash
pnpm dev
```

フロントエンドは http://localhost:5173 でアクセス可能になります

#### 開発用コマンド

```bash
# 開発サーバー起動
pnpm dev

# ビルド
pnpm build

# リント実行
pnpm lint

# API クライアント生成（バックエンドのSwagger更新後）
pnpm generate
```

### その他のサービス

#### MailHog (メール確認)
- Web UI: http://localhost:8025
- SMTP: localhost:1025

#### PostgreSQL
- メインDB: localhost:5432
- テストDB: localhost:15432
