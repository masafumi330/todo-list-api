# Tasks
- [x] POST /register
- [] ユーザーが新規登録したとき、Set-CookieでセッションIDを発行する
- [] POST /login
- [] ユーザーがログインしたとき、Set-CookieでセッションIDを発行する
- [] POST /logout
- [] ユーザーがログアウトしたとき、セッションIDを無効

## セッション管理
スキーマ定義
storage: redis
table: sessions
  - session_id: string (primary key)
  - user_id: integer (foreign key to users table)
  - created_at: datetime
  - expires_at: datetime
セッションの削除タイミング
  - 有効期限切れ
  - ユーザーログアウト時

# Goals
The skills you will learn from this project include:

- [] User authentication
- [] Schema design and Databases
- [] RESTful API design
- [] CRUD operations
- [] Error handling
- [] Security

# Requirements
You are required to develop a RESTful API with following endpoints

- [] User registration to create a new user
- [] Login endpoint to authenticate the user and generate a token
- [] CRUD operations for managing the to-do list
- [] Implement user authentication to allow only authorized users to access the to-do list
- [] Implement error handling and security measures
- [] Use a database to store the user and to-do list data (you can use any database of your choice)
- [] Implement proper data validation
- [] Implement pagination and filtering for the to-do list

# Spec
## Stack
- frontend: SPA(React)
- backend: Go
- Database: MySQL

## 画面
1. Top
2. ユーザー新規登録画面
3. ユーザーログイン画面
4. ToDo一覧画面

## API
1. ユーザー新規登録
```
POST /register
{
  "name": "John Doe",
  "email": "john@doe.com",
  "password": "password"
}

```

2. ユーザーログイン
```
POST /login
{
  "email": "john@doe.com",
  "password": "password"
}

```

3. ToDo作成
```
POST /todos
{
  "title": "Buy groceries",
  "description": "Buy milk, eggs, and bread"
}

```

4. ToDo更新
```

PUT /todos/1
{
  "title": "Buy groceries",
  "description": "Buy milk, eggs, bread, and cheese"
}
```

5. ToDo削除
```
DELETE /todos/1
```

6. Todoリスト取得
```
GET /todos?page=1&limit=10
{
  "data": [
    {
      "id": 1,
      "title": "Buy groceries",
      "description": "Buy milk, eggs, bread"
    },
    {
      "id": 2,
      "title": "Pay bills",
      "description": "Pay electricity and water bills"
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 2
}
```
