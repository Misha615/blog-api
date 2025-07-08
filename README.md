# Blog API

REST API для блога на Go (MySQL)

## 🚀 Запуск

1. Скопируйте `.env.example` в `.env` и настройте:
```ini
DB_USER=root
DB_PASS=password
DB_HOST=localhost:3306
DB_NAME=blog
```

2. Запустите:
```bash
go run cmd/main.go
```

## 📚 API Endpoints

| Метод | Путь              | Описание          |
|-------|-------------------|-------------------|
| GET   | /articles         | Все статьи        |
| POST  | /articles         | Создать статью    |
| PUT   | /articles/{id}    | Обновить статью   |
| DELETE| /articles/{id}    | Удалить статью    |
| GET   | /articles/{id}    | Выбрать статью    |

## 🛠 Технологии
- Go 1.21+
- MySQL
- Gorilla Mux