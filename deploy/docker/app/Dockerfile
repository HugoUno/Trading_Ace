FROM golang:1.20-alpine

WORKDIR /app

# 安裝基本工具
RUN apk add --no-cache git

# 複製 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源碼
COPY . .

# 編譯
RUN go build -o main cmd/server/main.go

# 執行
CMD ["./main"] 