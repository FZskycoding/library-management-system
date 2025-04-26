# 圖書館管理系統 Library Management System

這是一個使用 Go 語言開發的圖書館管理系統 API，使用PostgreSQL作為資料庫，透過 API 進行新增、查詢、刪除、更新書籍的動作，並透過 JWT 認證實現使用者登入、登出系統。

````

## 主要功能

1. 使用者管理
   - 使用者註冊
   - 使用者登入
   - 使用者登出
   - JWT 認證

2. 書籍管理
   - 新增書籍
   - 查詢書籍
   - 更新書籍資訊
   - 刪除書籍
   - 借書
   - 還書

## 技術棧

- 程式語言：Go
- 資料庫：PostgreSQL
- ORM：GORM
- 認證：JWT (JSON Web Token)

## 快速開始

1. 克隆專案
```bash
git clone https://github.com/FZskycoding/library-management-system
cd library-management-system
````

2. 設定資料庫

- 在 PostgreSQL 中創建資料庫
- 修改 `config/config.go` 中的資料庫連接資訊

3. 安裝依賴

```bash
go mod download
```

4. 運行專案

```bash
go run main.go
```


#### 註冊新用戶

```http
POST /api/auth/register
Content-Type: application/json

{
    "username": "user1",
    "password": "password123"
}
```

#### 用戶登入

```http
POST /api/auth/login
Content-Type: application/json

{
    "username": "user1",
    "password": "password123"
}
```

#### 新增書籍

```http
POST /api/books
Authorization: Bearer [token]
Content-Type: application/json

{
    "title": "Book Title",
    "author": "Author Name",
    "isbn": "1234567890"
}
```

