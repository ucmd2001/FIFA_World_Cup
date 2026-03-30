# FIFA World Cup Betting System

這是一個完整的前後端分離專案，旨在提供一個有趣的世足賽下注模擬平台。

## 目錄結構

- `/backend`: 基於 Go (Gin/GORM) 的 DDD 架構後端。
- `/frontend`: 基於 Vue 3 (Vite/Pinia) 的現代化前端。

## 快速啟動 (使用 Docker Compose)

這是最簡單的啟動方式，將會同時啟動後端、資料庫與前端。

```bash
docker-compose up --build
```

啟動後：
- 前端訪問位址：`http://localhost:3000`
- 後端 API 位址：`http://localhost:8080/api`

## 開發者設定

### 後端設定
請參考 `backend/README.md`。

### 前端設定
請參考 `frontend/README.md`。

## 功能預覽

1. **使用者註冊與登入**: 安全的 JWT 驗證。
2. **每日領取點數**: 獲取初始資金參與下注。
3. **賽事即時下注**: 動態計算賠率，體驗真實下注感。
4. **管理員控制台**: 賽事結算、使用者權限管理。
