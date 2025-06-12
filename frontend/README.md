# Frontend

農地・農業用施設等災害復旧支援システムのフロントエンド

React + TypeScript + ViteによるモダンなWebアプリケーション

## 技術スタック

- **React**: ユーザーインターフェース構築のためのJavaScriptライブラリ
- **TypeScript**: スケーラブルな型付きJavaScript
- **Vite**: 次世代フロントエンド開発ツール
- **Tailwind CSS**: ユーティリティファーストCSSフレームワーク
- **shadcn/ui**: Radix UIとTailwind CSSで構築された美しくデザインされたコンポーネント
- **TanStack Query**: データフェッチ、キャッシュ、更新のための強力な非同期状態管理
- **TanStack Router**: Reactアプリケーション向け型安全ルーティング
- **Orval**: OpenAPI仕様からAPIクライアント自動生成
- **Biome**: 高速で信頼性の高いリンター・フォーマッター
- **Zod**: 静的型推論を持つTypeScriptファーストスキーマ検証
- **Zustand**: 小さく、高速でスケーラブルな状態管理ソリューション

## プロジェクト構成

プロジェクトは機能ベースアーキテクチャに従っています：

```
src/
├── api/                  # API関連コード
│   └── mutator/          # カスタムAPIクライアント設定
├── components/           # 共有コンポーネント
│   └── ui/               # UIコンポーネント（shadcn/ui）
├── features/             # 機能モジュール
│   └── auth/             # 認証機能
│       ├── components/   # 機能固有コンポーネント
│       ├── hooks/        # 機能固有フック
│       └── types/        # 機能固有型定義
├── lib/                  # ユーティリティ関数と共有コード
│   ├── validations/      # Zodスキーマ
│   └── utils.ts          # ユーティリティ関数
├── routes/               # アプリケーションルート
├── store/                # グローバル状態管理
├── App.tsx               # メインアプリケーションコンポーネント
└── main.tsx              # アプリケーションエントリーポイント
```

## はじめに

### 前提条件

- Node.js（v16以上）
- npm または yarn

### インストール

1. リポジトリをクローン
2. 依存関係をインストール：

```bash
cd frontend
npm install
# または
yarn
```

### 開発

開発サーバーを起動：

```bash
npm run dev
# または
yarn dev
```

### 本番環境用ビルド

本番環境用にアプリケーションをビルド：

```bash
npm run build
# または
yarn build
```

### リント・フォーマット

コードのリント：

```bash
npm run lint
# または
yarn lint
```

コードのフォーマット：

```bash
npm run format
# または
yarn format
```

## 新機能の追加

新機能を追加するには：

1. `src/features/`に新しいディレクトリを作成
2. 機能固有のコンポーネント、フック、型定義を追加
3. 必要に応じて`src/routes/`にルートを作成
4. アプリケーションで機能をインポート・使用

## API統合

1. OpenAPI仕様書をルートディレクトリの`openapi.yaml`に配置
2. OrvalでAPIクライアントを生成：

```bash
npm run generate-api
# または
yarn generate-api
```

## コンポーネントのスタイリング

プロジェクトではスタイリングにTailwind CSSを使用しています。複雑なコンポーネントの場合は、shadcn/uiの使用を検討するか、`components/ui`ディレクトリにカスタムコンポーネントを作成してください。

## 状態管理

- グローバル状態管理にはZustandを使用
- サーバー状態にはTanStack Queryを使用
- コンポーネントローカル状態にはReactの組み込み状態管理（useState、useReducer）を使用

## ディレクトリ詳細

### `src/api/`
- **mutator/**: APIクライアントのカスタム設定
- 自動生成されたAPIクライアントコード
- エラーハンドリングとインターセプター設定

### `src/components/`
- **ui/**: 再利用可能なUIコンポーネント
- ボタン、モーダル、フォーム要素などの基本コンポーネント
- shadcn/uiコンポーネントライブラリ

### `src/features/`
各機能モジュールの構成例：
```
features/
├── auth/                 # 認証機能
│   ├── components/       # ログインフォーム、認証ガードなど
│   ├── hooks/            # useAuth、useLoginなど
│   └── types/            # User、LoginRequest型など
├── disaster/             # 被災情報管理
│   ├── components/       # 被災情報フォーム、一覧表示など
│   ├── hooks/            # useDisasters、useCreateDisasterなど
│   └── types/            # Disaster、DisasterForm型など
└── dashboard/            # ダッシュボード
    ├── components/       # 統計カード、チャートなど
    ├── hooks/            # useDashboardData
    └── types/            # DashboardStats型など
```

### `src/lib/`
- **validations/**: Zodを使用したスキーマ検証
- **utils.ts**: 共通ユーティリティ関数
- API設定、日付フォーマット、型ガードなど

### `src/routes/`
TanStack Routerを使用したルート定義：
```
routes/
├── __root.tsx            # ルートレイアウト
├── index.tsx             # ホームページ
├── login.tsx             # ログインページ
├── dashboard/            # ダッシュボード関連ルート
├── disasters/            # 被災情報関連ルート
└── settings/             # 設定関連ルート
```

### `src/store/`
Zustandを使用したグローバル状態管理：
```
store/
├── authStore.ts          # 認証状態
├── uiStore.ts            # UI状態（モーダル、サイドバーなど）
└── index.ts              # ストア統合
```

## 開発ガイドライン

### コンポーネント作成
```typescript
// 関数コンポーネントを使用
const MyComponent: React.FC<Props> = ({ prop1, prop2 }) => {
  return <div>{/* JSX */}</div>;
};

// propsの型定義を明確に
interface Props {
  title: string;
  onClick: () => void;
}
```

### カスタムフック作成
```typescript
// useから始まる命名
const useMyFeature = () => {
  const [state, setState] = useState();
  
  return {
    state,
    actions: {
      updateState: setState
    }
  };
};
```

### API呼び出し
```typescript
// TanStack Queryを使用
const { data, isLoading, error } = useQuery({
  queryKey: ['disasters'],
  queryFn: () => api.getDisasters()
});
```

## パフォーマンス最適化

### コード分割
```typescript
// React.lazy を使用した動的インポート
const Dashboard = lazy(() => import('./features/dashboard/Dashboard'));
```

### メモ化
```typescript
// React.memo でコンポーネントをメモ化
const ExpensiveComponent = memo(({ data }) => {
  // 重い処理...
});

// useMemo でデータをメモ化
const processedData = useMemo(() => {
  return expensiveCalculation(data);
}, [data]);
```

## テスト戦略

### ユニットテスト
```typescript
// components/__tests__/MyComponent.test.tsx
describe('MyComponent', () => {
  it('renders correctly', () => {
    render(<MyComponent />);
    expect(screen.getByText('Expected Text')).toBeInTheDocument();
  });
});
```

### 統合テスト
```typescript
// features/__tests__/auth.integration.test.tsx
describe('Authentication Flow', () => {
  it('allows user to login', async () => {
    // ログインフローのテスト
  });
});
```

## デプロイ

### 環境設定
```bash
# 本番環境用環境変数
VITE_API_URL=https://api.example.com
VITE_APP_ENV=production
```

### ビルド最適化
```typescript
// vite.config.ts
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],
          ui: ['@radix-ui/react-dialog', '@radix-ui/react-dropdown-menu']
        }
      }
    }
  }
});
```

## 利用可能なnpmスクリプト

```bash
# 開発
npm run dev              # 開発サーバー起動
npm run build            # 本番ビルド
npm run preview          # ビルド結果プレビュー

# コード品質
npm run lint             # リント実行
npm run lint:fix         # リント自動修正
npm run format           # コードフォーマット
npm run type-check       # TypeScript型チェック

# テスト
npm run test             # テスト実行
npm run test:watch       # テスト監視モード
npm run test:coverage    # カバレッジ付きテスト

# API
npm run generate-api     # OpenAPIからクライアント生成
npm run validate-api     # API仕様書検証
```