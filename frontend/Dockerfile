# --- ビルドステージ ---
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# --- 実行ステージ ---
FROM nginx:stable-alpine

COPY --from=builder /app/dist /usr/share/nginx/html

# 必要であれば nginx の設定ファイルをコピー (PoCであれば省略可能)
# COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]