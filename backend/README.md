# FIFA World Cup Betting System - Backend

這是一個基於 Go 語言開發的世足賽下注系統後端，採用領域驅動設計 (Domain-Driven Design, DDD) 概念進行架構開發。

## 技術棧

- **語言**: Go 1.22+
- **框架**: Gin Web Framework
- **資料庫**: SQLite (使用 GORM ORM)
- **認證**: JWT (JSON Web Token)
- **定時任務**: Robfig/Cron (用於賽事同步與加分)

## 目錄結構 (DDD)

```text
internal/
├── auth/           # 認證與授權 (Login, Register, JWT Middleware)
├── user/           # 使用者管理 (Profile, Leaderboard, Points)
├── match/          # 賽事管理 (Sync, Resolve, Pool Stats)
├── bet/            # 下注邏輯 (Place Bet, History)
├── common/         # 共用組件 (Response, Errors, Utils)
├── infrastructure/ # 基礎設施 (DB, Scheduler)
└── router/         # 路由配置
```

## 快速開始

### 1. 環境變數配置
複製 `.env.example` 並重新命名為 `.env`，填入必要的資訊：
- `PORT`: 後端監聽埠 (預設 8080)
- `DB_PATH`: SQLite 檔案路徑
- `JWT_SECRET`: JWT 簽名金鑰
- `FOOTBALL_DATA_API_KEY`: football-data.org 的 API Key (可選)

### 2. 本地執行
```bash
go mod download
go run cmd/server/main.go
```

## API 設計規範

- **統一回應格式**: 所有 API 回傳皆符合以下結構：
  ```json
  {
    "status": 200,
    "message": "Success",
    "data": { ... }
  }
  ```
- **身份驗證**: 需在 Header 帶入 `Authorization: Bearer <TOKEN>`。

## 重點功能

- **賽事自動同步**: 定時抓取外部 API 更新賽程。
- **動態賠率**: 下注金額池即時計算，根據贏家比例分配輸家注金。
- **每日領取**: 使用者每日可領取 1000 點數參與下注。
