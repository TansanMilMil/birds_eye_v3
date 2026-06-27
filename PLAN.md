# birds_eye_v3 モノレポ移行プラン

## 概要

`birds-eye-page`（FE）と `birds_eye_v3`（BE）を単一リポジトリに統合する。

---

## 決定事項

### 1. ディレクトリ構造

フラット分割（案A）を採用。

```
birds_eye_v3/
├── frontend/          # birds-eye-page のコード一式
│   ├── src/
│   ├── index.html
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── package.json
│   └── Dockerfile
├── backend/           # birds_eye_v3 のコード一式
│   ├── go/
│   │   └── src/
│   ├── nginx/
│   ├── mysql/
│   ├── scripts/
│   ├── go.mod
│   ├── go.sum
│   └── go-entrypoint.sh
├── docker-compose.yml
├── .env
├── .env.example
└── README.md
```

**移行対象外:** loki / promtail / grafana（監視系は既存リポジトリに残す）

---

### 2. docker-compose 構成

**単一の `docker-compose.yml` + `.env`** で dev/prod 両方を賄う。

dev/prod の切り替えは `BIRDSEYE_BIRDSEYEAPI_EXECUTION_MODE` 等の環境変数で制御する（既存の仕組みを踏襲）。override ファイルや profiles は使わない。

サービス構成:

| サービス | 役割 |
|---|---|
| `frontend` | `vite build --watch` でビルドし続ける（dev/prod共通） |
| `nginx` | 静的ファイル配信 + API プロキシ |
| `go` | Gin APIサーバー |
| `mysql` | データベース |
| `selenium` | スクレイピング用ブラウザ |

---

### 3. フロントエンドの配信方式

**dev/prod ともに nginx から静的ファイルを配信する**（CloudFront / S3 デプロイは廃止）。

- `frontend` サービスが `vite build --watch` を実行し、`build/` を Docker volume に出力し続ける
- `nginx` サービスがその volume をマウントして静的ファイルを配信する
- nginx は `/api/` へのリクエストを `go:8080` にプロキシする

```
[ブラウザ]
   └→ nginx
       ├ /        → frontend の build/ （静的ファイル）
       └ /api/    → go:8080 （APIプロキシ）
```

---

### 4. .env 戦略

ルートに**単一の `.env`** を置き、FE・BE 両方の変数をまとめて管理する。

```env
# --- Frontend ---
BIRDSEYE_VITE_BIRDS_EYE_API_ENDPOINT=http://localhost/api

# --- Backend ---
BIRDSEYE_MYSQL_ROOT_PASSWORD=
BIRDSEYE_GO_API_PORT=8080
BIRDSEYE_BIRDSEYEAPI_EXECUTION_MODE=dev
BIRDSEYE_OPENAI_MODEL=gpt-4-turbo
BIRDSEYE_OPENAI_CHAT_ENDPOINT=https://api.openai.com/v1/chat/completions
BIRDSEYE_BIRDSEYEAPI_V2_OPENAI_API_KEY=
BIRDSEYE_BIRDSEYEAPI_V2_CLAUDE_API_KEY=
BIRDSEYE_AWS_REGION=
BIRDSEYE_AWS_CLOUDFRONT_BIRDSEYEAPIPROXY_DISTRIBUTION_ID=
BIRDSEYE_SCRAPING_ARTICLES=10
BIRDSEYE_SELENIUM_URL=http://selenium:4444/wd/hub
```

---

## 移行対象外

- loki / promtail / grafana（既存 birds_eye_v3 リポジトリで引き続き管理）
- CI/CD（GitHub Actions）は移行プランに含めず、別途検討する
