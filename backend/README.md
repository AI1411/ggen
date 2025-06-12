## ディレクトリ構成

```
backend/
├── cmd/                          # エントリーポイント（実行可能なファイル）
│   ├── api/                      # APIサーバーのmain
│   │   └── main.go              # APIサーバー起動ファイル
│   ├── gormgen/                 # ORM自動生成コマンド
│   │   ├── generate_all/        # 全モデル生成
│   │   └── generate_associations/ # アソシエーション生成
│   └── seed/                    # データ投入コマンド
│       └── municipality/        # 自治体データ投入
├── docs/                        # APIドキュメント
│   └── api/                     # Swaggerドキュメント
├── internal/                    # 内部パッケージ（非公開）
│   ├── di/                      # 依存性注入
│   │   └── provider.go          # DIコンテナ設定
│   ├── domain/                  # ドメイン層
│   │   ├── model/               # ドメインモデル（GORM自動生成）
│   │   ├── query/               # クエリビルダー（GORM自動生成）
│   │   └── repository/          # リポジトリインターフェース
│   ├── env/                     # 環境変数管理
│   ├── errors/                  # エラー処理
│   │   └── error.go             # カスタムエラー定義
│   ├── handler/                 # プレゼンテーション層（HTTPハンドラー）
│   │   ├── common_handler.go    # 共通ハンドラー
│   │   ├── error_response.go    # エラーレスポンス
│   │   ├── prefecture_handler.go # 都道府県関連API
│   │   └── validator.go         # バリデーション
│   ├── infra/                   # インフラストラクチャ層
│   │   ├── datastore/           # データベース実装
│   │   ├── db/                  # データベース接続
│   │   └── logger/              # ログ出力
│   ├── server/                  # サーバー設定
│   │   ├── middleware/          # ミドルウェア
│   │   └── route.go             # ルーティング設定
│   └── usecase/                 # アプリケーション層（ユースケース）
├── migrations/                  # データベースマイグレーション
├── tests/                       # テスト関連
│   ├── mock/                    # モックファイル
│   └── testutils/               # テストユーティリティ
├── go.mod                       # Goモジュール定義
└── go.sum                       # Goモジュール依存関係
```

## アーキテクチャ

このプロジェクトは **クリーンアーキテクチャ** を採用しており、以下の層に分かれています：

### 1. プレゼンテーション層 (`handler/`)

- HTTPリクエストの受信・レスポンス返却
- リクエストのバリデーション
- レスポンスの整形

### 2. アプリケーション層 (`usecase/`)

- ビジネスロジックの実装
- ドメインオブジェクトの操作
- トランザクション制御

### 3. ドメイン層 (`domain/`)

- **model/**: データベーステーブルに対応するモデル（GORM自動生成）
- **query/**: 型安全なクエリビルダー（GORM自動生成）
- **repository/**: データアクセスのインターフェース定義

### 4. インフラストラクチャ層 (`infra/`)

- **datastore/**: リポジトリインターフェースの実装
- **db/**: データベース接続管理
- **logger/**: ログ出力実装

## 重要なファイル・ディレクトリの詳細

### `cmd/`

- **api/main.go**: APIサーバーのエントリーポイント
- **gormgen/**: GORMのコード自動生成ツール
- **seed/**: 初期データ投入ツール

### `internal/di/`

- 依存性注入（Dependency Injection）の設定
- 各層の依存関係を管理

### `internal/domain/`

- **model/**: データベーステーブルと1:1対応するGORM構造体
- **query/**: 型安全なクエリビルダー（GORMのGen機能）
- **repository/**: データアクセス層のインターフェース

### `internal/handler/`

- RESTful APIのハンドラー実装
- リクエスト・レスポンスの処理
- HTTPステータスコードの管理

### `internal/usecase/`

- ビジネスロジックの中核
- 複数のリポジトリを組み合わせた処理
- アプリケーション固有のルール実装

### `internal/infra/`

- 外部システムとの接続実装
- データベース操作の具体的な実装
- ログ出力やその他のインフラ機能

### `migrations/`

- データベーススキーマのバージョン管理
- SQLファイルによるマイグレーション定義

### `tests/`

- **mock/**: テスト用のモックオブジェクト
- **testutils/**: テスト共通機能

## 開発フロー

1. **マイグレーション**: `migrations/`でスキーマ定義
2. **モデル生成**: `make generate-models`でGORMモデル自動生成
3. **リポジトリ**: `domain/repository/`でインターフェース定義
4. **実装**: `infra/datastore/`でリポジトリ実装
5. **ユースケース**: `usecase/`でビジネスロジック実装
6. **ハンドラー**: `handler/`でHTTP API実装
7. **テスト**: 各層でユニットテスト作成

## 技術スタック

### フレームワーク・ライブラリ
- **Webフレームワーク**: Gin
- **ORM**: GORM v2
- **依存性注入**: fx
- **バリデーション**: go-playground/validator
- **ログ**: slog

### データベース
- **プライマリDB**: PostgreSQL
- **マイグレーション**: golang-migrate

### テスト
- **テストフレームワーク**: 標準testing + testify
- **モック生成**: GoMock

### ドキュメント
- **API仕様書**: Swagger/OpenAPI 2.0
- **コード生成**: swaggo

## セットアップ手順

1. **環境準備**
```bash
# 依存関係インストール
go mod download

# データベース起動（Docker Compose）
docker-compose up -d db
```

2. **データベースセットアップ**
```bash
# マイグレーション実行
make migrate-up

# シードデータ投入
make seed
```

3. **コード生成**
```bash
# GORMモデル生成
make generate-models

# モック生成
make generate-mocks

# Swagger文書生成
make generate-docs
```

4. **サーバー起動**
```bash
# 開発サーバー起動
make run

# または
go run cmd/api/main.go
```

## 利用可能なMakeコマンド

```bash
# 開発用
make run                 # サーバー起動
make test               # テスト実行
make lint               # リント実行

# コード生成
make generate-models    # GORMモデル生成
make generate-mocks     # モック生成
make generate-docs      # Swagger文書生成

# データベース
make migrate-up         # マイグレーション実行
make migrate-down       # マイグレーションロールバック
make seed              # シードデータ投入

# ビルド・デプロイ
make build             # バイナリビルド
make docker-build      # Dockerイメージビルド
```

## API エンドポイント

### 認証
- `POST /api/auth/login` - ログイン
- `POST /api/auth/logout` - ログアウト

### 都道府県管理
- `GET /api/prefectures` - 都道府県一覧取得
- `GET /api/prefectures/{code}` - 都道府県詳細取得

### 被災情報管理
- `GET /api/disasters` - 被災情報一覧取得
- `POST /api/disasters` - 被災情報登録
- `GET /api/disasters/{id}` - 被災情報詳細取得
- `PUT /api/disasters/{id}` - 被災情報更新
- `DELETE /api/disasters/{id}` - 被災情報削除

### 写真管理
- `POST /api/disasters/{id}/photos` - 写真アップロード
- `GET /api/disasters/{id}/photos` - 写真一覧取得

詳細なAPI仕様書は `http://localhost:8080/swagger/` で確認できます。
