# TT-MEDODS
Technical Task Backend Developer

## Auth Service implemented on Golang

Develop on Platform: Ubuntu 24.04 (WSL)  
Go version: `1.23`  

## Project startup
1. Настраиваем env:
```
cp .env.example .env
```
2. Настраиваем env
3. Запускаем:
```
docker compose up
```

В режиме разработки:
```
docker compose up --build
or
docker compose build --no-cache
docker compose up
```

## Endpoints
- Service host: `localhost:8080`  
- ID тестового пользователя: `5f32a8a2-45f4-4c0c-9f7e-61e8b8d4edc9` 
- Так же приложен postman 


### 1. Issue Tokens
**Endpoint**: `GET /auth/token/{user_id}`

**Description**: 
Получение новых токенов для пользователя по его `user_id`.

**Parameters**:
- `user_id` (Path): Уникальный идентификатор пользователя.

**Response**:
- Возвращает JSON обьект пары токенов `access_token` and `refresh_token`.

### 2. Refresh Tokens
**Endpoint**: `PATCH /auth/token/{user_id}`

**Description**: 
Обновление существующих токенов. Необходимо передать как `access_token`, так и `refresh_token` для успешного обновления.

**Body**:
```
{
  "access_token": "string",
  "refresh_token": "string"
}
```

**Response**:
- Возвращает JSON обьект пары токенов `access_token` and `refresh_token`.


### P.S.
> Изначалено думал реализовывать на fiber, но в итоге реализовал на http через mux.  
> JWT реализовывал по стандарту: RFC 7519
> Middleware не накидывал специально  
> Проверял на удалённом сервере бд

2024 (c) statebyte