# birds_eye_v3

IT ニュースを自動収集・AI 要約して閲覧できる Web アプリ。  
フロントエンド・バックエンド・インフラ設定を単一リポジトリで管理するモノレポ構成。

---

## 機能

- 複数ソースから IT ニュースを自動スクレイピング
- OpenAI / Claude による日本語 200 字要約
- Hatena Bookmark の反応数・コメント取得（Selenium）
- Google Trends のトレンドワード表示
- ダーク / ライトテーマ切り替え

---

## アーキテクチャ

### 本番環境

```
ユーザー
  │
  ▼
┌──────────────────┐
│   CloudFront     │  CDN・HTTPS 終端・キャッシュ
└────────┬─────────┘
         │ HTTP (origin)
         ▼
┌─────────────────────────────────────────────┐
│              さくら VPS                      │
│                                             │
│  ┌──────────────────────────────────────┐   │
│  │           Docker Compose             │   │
│  │                                      │   │
│  │  ┌────────────┐                      │   │
│  │  │   Nginx    │ :8082                │   │
│  │  │            │                      │   │
│  │  │  /         │──→ frontend-build    │   │
│  │  │  /api/     │──→ go:8080           │   │
│  │  │  /Health   │──→ 200 ok (内部処理) │   │
│  │  └────────────┘                      │   │
│  │         │                            │   │
│  │  ┌──────┴──────┐  ┌───────────────┐  │   │
│  │  │  Go (Gin)   │  │   Frontend    │  │   │
│  │  │   :8080     │  │ vite build    │  │   │
│  │  └──────┬──────┘  │   --watch     │  │   │
│  │         │         └───────┬───────┘  │   │
│  │  ┌──────┴──────┐          │          │   │
│  │  │    MySQL    │  frontend-build     │   │
│  │  │    :3306    │  (Docker volume)    │   │
│  │  └─────────────┘                    │   │
│  │                                      │   │
│  │  ┌─────────────┐                     │   │
│  │  │  Selenium   │ :4444 (内部のみ)    │   │
│  │  │  Firefox    │                     │   │
│  │  └─────────────┘                     │   │
│  └──────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
```

**ポイント:**
- CloudFront がHTTPS を終端し、オリジンの Nginx へ HTTP で転送する
- Nginx が静的ファイル配信と API プロキシを一手に担う
- フロントエンドコンテナは `vite build --watch` でビルドし続け、成果物を Docker volume 経由で Nginx に渡す
- Selenium は外部に公開せず、Go API からのみアクセスする

### リクエストフロー

```
GET /           → CloudFront → Nginx → frontend-build volume (React SPA)
GET /api/news/* → CloudFront → Nginx → Go:8080 → MySQL
POST /api/news/scrape → Go → Selenium(Firefox) → 各ニュースサイト
```

---

## コンポーネント一覧

| サービス | イメージ | ホストポート | 役割 |
|---|---|---|---|
| nginx | nginx:1.31.2 | 8082 | リバースプロキシ・静的ファイル配信 |
| go | golang:1.24 | 8080 | Gin API サーバー |
| mysql | mysql:9.3 | 3307 | データ永続化 |
| selenium | selenium/standalone-firefox:133.0 | 4444 | ヘッドレス Firefox（スクレイピング用） |
| frontend | (カスタム) | — | Vite ビルドワーカー（静的ファイル生成のみ） |

---

## ディレクトリ構成

```
birds_eye_v3/
├── frontend/              # React (TypeScript) + Vite + MUI
│   ├── src/
│   │   ├── pages/         # /news, /trends ページ
│   │   ├── api/           # API クライアント
│   │   ├── store/         # Redux store
│   │   └── share-components/
│   ├── Dockerfile
│   └── vite.config.ts
├── backend/               # Go + Gin API サーバー
│   ├── go/src/            # アプリケーションコード
│   ├── go/dist/           # ビルド済みバイナリ（本番用）
│   ├── build.sh           # バイナリビルドスクリプト
│   ├── go-entrypoint.sh   # コンテナ起動スクリプト
│   └── scrape.sh          # リモートスクレイプトリガー
├── nginx/
│   ├── nginx.conf
│   └── conf.d/default.conf
├── mysql/
│   └── create_db.sql      # DB 初期化 SQL
├── docker-compose.yml
├── .env                   # 環境変数（git 管理外）
└── .env.example           # 環境変数テンプレート
```

---

## セットアップ

### 前提条件

- Docker / Docker Compose
- OpenAI API キー または Anthropic Claude API キー（記事要約に使用）
- AWS アカウント（CloudFront キャッシュ無効化に使用）

### 手順

```bash
# 1. リポジトリをクローン
git clone <repo-url> birds_eye_v3
cd birds_eye_v3

# 2. 環境変数ファイルを作成
cp .env.example .env
# .env を編集して各値を設定

# 3. コンテナを起動
docker compose up -d

# 4. ログ確認
docker compose logs -f
```

起動後、`http://localhost:8082` でアクセスできる。

---

## 環境変数

`.env.example` を参照。主要な変数は以下のとおり。

| 変数名 | 説明 |
|---|---|
| `MYSQL_ROOT_PASSWORD` | MySQL root パスワード |
| `BIRDSEYEAPI_V2_OPENAI_API_KEY` | OpenAI API キー（要約用） |
| `BIRDSEYEAPI_V2_CLAUDE_API_KEY` | Claude API キー（要約用） |
| `AWS_REGION` | AWS リージョン |
| `AWS_CLOUDFRONT_BIRDSEYEAPIPROXY_DISTRIBUTION_ID` | CloudFront ディストリビューション ID |
| `BIRDSEYEAPI_EXECUTION_MODE` | `PRODUCTION` で本番バイナリを実行、それ以外は開発シェル |
| `VITE_BIRDS_EYE_API_ENDPOINT` | フロントエンドから叩く API エンドポイント URL |
| `SCRAPING_ARTICLES` | スクレイプする記事数（デフォルト 10） |

---

## 本番デプロイ

### バイナリビルド

```bash
cd backend
./build.sh
# → go/dist/birds_eye_v3 が生成される
docker compose restart go
```

### スクレイプ手動トリガー

```bash
cd backend
./scrape.sh
# VENUS_SSH_HOST が .env に設定されている必要がある
```

---

## API エンドポイント

| メソッド | パス | 説明 |
|---|---|---|
| GET | `/api/news/:target_date` | 指定日付（YYYY-MM-DD）のニュース一覧 |
| GET | `/api/news/news-reactions/:news-id` | 指定記事の Hatena Bookmark 反応 |
| POST | `/api/news/scrape` | ニュース・反応のスクレイプ実行（多重実行は 409） |
| GET | `/api/news/trends` | Google Trends のトレンドワード |
| GET | `/HealthCheck` | Nginx が直接 200 を返すヘルスチェック |

---

## データソース

- [CloudWatch by Impress](https://cloud.watch.impress.co.jp/)
- [Hatena Bookmark（IT 注目エントリー）](https://b.hatena.ne.jp/hotentry/it)
- [Zenn](https://zenn.dev/)
- [ZDNet Japan](https://japan.zdnet.com/)

各ソースから最大 15 件を取得し、OpenAI `gpt-4.1-mini` で日本語 200 字に要約する。

---

## 技術スタック

| レイヤー | 技術 |
|---|---|
| フロントエンド | React 18 / TypeScript / Vite / MUI / Redux Toolkit |
| バックエンド | Go 1.24 / Gin |
| データベース | MySQL 9.3 |
| スクレイピング | Selenium (headless Firefox) |
| リバースプロキシ | Nginx 1.31 |
| CDN | AWS CloudFront |
| インフラ | さくら VPS + Docker Compose |
| AI 要約 | OpenAI gpt-4.1-mini / Anthropic claude-3-5-sonnet |
