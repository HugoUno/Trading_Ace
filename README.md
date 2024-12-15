# Trading Ace

這是一個追蹤 Uniswap V2 用戶交易的積分系統，讓用戶可以透過交易賺取積分。

## 這系統在幹嘛？

### 新手任務
- 只要在 WETH/USDC 池交易超過 1000u 就可以拿到 100 積分
- 一次就好，後面不會再給了喔

### 分享池任務
- 每週會看你在 WETH/USDC 池的交易量佔比
- 根據佔比分享 10,000 積分
- 要先完成新手任務才能拿這個獎勵喔

### 目標交易池
- 在這個池子交易才算：WETH/USDC
- 合約位址：0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc

## 怎麼跑起來？

### 1. 先準備這些東西
- Go 1.20 以上
- Docker 跟 Docker Compose
- 一個 Infura API Key（要連以太坊用的）

### 2. 安裝跟設定

```bash
# 裝一下需要的套件
go mod init trading-ace
go mod tidy

# 複製設定檔，然後改一下裡面的設定
cp config.yaml.example config.yaml
vim config.yaml  # 記得填入你的 Infura API Key
```

### 3. 啟動服務

```bash
# 先把資料庫跟 Redis 跑起來
docker compose up -d

# 等它們都準備好了再跑主程式
go run cmd/server/main.go
```

## API 怎麼用？

### 查詢用戶的任務狀態

```http
GET /api/v1/users/{address}/tasks

# 回傳範例：
{
    "tasks": [
        {
            "id": 1,
            "type": "ONBOARDING",
            "status": "COMPLETED",
            "points": 100,
            "amount": 1500.50
        }
    ]
}
```

### 查積分紀錄

```http
GET /api/v1/users/{address}/points/history

# 回傳範例：
{
    "history": [
        {
            "task_id": 1,
            "points": 100,
            "timestamp": "2024-03-01T10:00:00Z",
            "type": "ONBOARDING"
        }
    ]
}
```

### 看排行榜

```http
GET /api/v1/leaderboard

# 回傳範例：
{
    "rankings": [
        {
            "rank": 1,
            "address": "0x...",
            "total_points": 2600
        }
    ]
}
```

## 開發者專區

### 常用指令

```bash
# 跑測試
go test ./... -cover

# 編譯程式
go build -o main cmd/server/main.go

# 整理程式碼
go fmt ./...
```

### 環境變數

```env
# 資料庫設定
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=trading_ace
DB_HOST=localhost
DB_PORT=5432

# 以太坊設定
ETH_NODE_URL=https://mainnet.infura.io/v3/你的-PROJECT-ID
ETH_POOL_ADDRESS=0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc

# 伺服器設定
SERVER_PORT=8080
```

## 專案結構說明

```
trading-ace/
├── cmd/                    # 主程式
├── internal/              # 內部程式碼
│   ├── config/           # 設定檔
│   ├── domain/           # 資料結構
│   ├── repository/       # 資料庫存取
│   ├── service/          # 商業邏輯
│   └── handler/          # API 處理
├── pkg/                   # 共用程式
├── scripts/              # 腳本
└── deploy/               # 部署用檔案
```

