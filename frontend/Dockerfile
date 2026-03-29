# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

# 執行編譯（會產生 dist 目錄）
RUN npm run build

# Final stage - 使用 Nginx 託管靜態檔案
FROM nginx:alpine

# 複製編譯後的檔案至 Nginx 目錄
COPY --from=builder /app/dist /usr/share/nginx/html

# 複製 Nginx 設定 (自定義)
# COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
