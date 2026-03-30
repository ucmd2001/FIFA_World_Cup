# FIFA World Cup Betting System - Frontend

這是一個基於 Vue 3 開發的世足賽下注系統前端，提供直觀的即時賽況、下注與排行榜介面。

## 技術棧

- **框架**: Vue 3 (Composition API)
- **構建工具**: Vite
- **狀態管理**: Pinia
- **路由**: Vue Router
- **網路請求**: Axios
- **樣式**: Vanilla CSS (現代化深色/淺色主題支援)

## 主要功能

- **即時賽事列表**: 查看進行中、已結束與未開始的賽事。
- **下注介面**: 簡單易用的下注彈窗與餘額檢視。
- **個人紀錄**: 追蹤歷史下注紀錄與勝負狀態。
- **英雄榜**: 全站使用者積分排名。
- **後台管理**: 提供管理員新增賽事、結算賽果與管理使用者的功能。

## 快速開始

### 1. 環境變數配置
複製 `.env.example` 並重新命名為 `.env`：
- `VITE_API_BASE_URL`: 後端 API 的 URL (例如 `http://localhost:8080/api`)

### 2. 安裝與啟動
```bash
npm install
npm run dev
```

## 專案亮點

- **Responsive Design**: 完美支援桌面與行動裝置。
- **Dark Mode**: 預設支援深淺色模式切換，提升視覺體驗。
- **Micro-animations**: 平滑的過渡效果與操作回饋。
