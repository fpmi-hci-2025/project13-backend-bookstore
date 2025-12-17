# BookStore Backend API

REST API для книжного онлайн-магазина, разработанное на Go с использованием Gin framework.

## 🚀 Особенности

- **RESTful API** с полной документацией Swagger
- **JWT аутентификация** для защищённых эндпоинтов
- **PostgreSQL** база данных с миграциями
- **Docker** для локальной разработки и деплоя
- **CI/CD** с GitHub Actions
- **Проверка качества кода** с SonarCloud и golangci-lint

## 📋 API Endpoints

### Аутентификация
- `POST /api/v1/auth/register` - Регистрация пользователя
- `POST /api/v1/auth/login` - Вход в систему
- `GET /api/v1/me` - Получить профиль текущего пользователя

### Книги
- `GET /api/v1/books` - Список всех книг
- `GET /api/v1/books/:id` - Получить книгу по ID
- `GET /api/v1/books/search?q=query` - Поиск книг
- `POST /api/v1/books` - Создать книгу (требуется авторизация)
- `PUT /api/v1/books/:id` - Обновить книгу
- `DELETE /api/v1/books/:id` - Удалить книгу

### Авторы
- `GET /api/v1/authors` - Список авторов
- `GET /api/v1/authors/:id` - Получить автора с его книгами
- `POST /api/v1/authors` - Создать автора
- `PUT /api/v1/authors/:id` - Обновить автора
- `DELETE /api/v1/authors/:id` - Удалить автора

### Корзина
- `GET /api/v1/cart` - Получить корзину пользователя
- `POST /api/v1/cart` - Добавить книгу в корзину
- `PUT /api/v1/cart/:id` - Обновить количество
- `DELETE /api/v1/cart/:id` - Удалить из корзины
- `DELETE /api/v1/cart` - Очистить корзину

### Заказы
- `GET /api/v1/orders` - Список заказов пользователя
- `POST /api/v1/orders` - Создать заказ
- `GET /api/v1/orders/:id` - Получить заказ
- `POST /api/v1/orders/:id/cancel` - Отменить заказ

### Избранное
- `GET /api/v1/favorites` - Список избранных книг
- `POST /api/v1/favorites` - Добавить в избранное
- `DELETE /api/v1/favorites/:book_id` - Удалить из избранного

### Отзывы
- `GET /api/v1/books/:id/reviews` - Отзывы на книгу
- `POST /api/v1/reviews` - Создать отзыв
- `PUT /api/v1/reviews/:id` - Обновить отзыв
- `DELETE /api/v1/reviews/:id` - Удалить отзыв

## 🛠 Технологии

- **Go 1.23** - язык программирования
- **Gin** - веб-фреймворк
- **PostgreSQL 16** - база данных
- **pgx/v5** - драйвер PostgreSQL
- **JWT** - аутентификация
- **Swagger** - документация API
- **Docker** - контейнеризация
- **GitHub Actions** - CI/CD
- **SonarCloud** - анализ качества кода

## 🏃 Быстрый старт

### Требования
- Go 1.23+
- Docker и Docker Compose
- Make (опционально)

### Локальная разработка

1. **Клонировать репозиторий**
```bash
git clone https://github.com/fpmi-hci-2025/project13-backend-bookstore.git
cd project13-backend-bookstore
```

2. **Запустить базу данных**
```bash
make dev-db
# или
docker compose -f docker-compose.dev.yml up -d
```

3. **Настроить переменные окружения**
```bash
cp .env.example .env
# Отредактировать .env при необходимости
```

4. **Запустить приложение**
```bash
make run
# или
go run ./cmd/api
```

5. **Открыть Swagger документацию**
```
http://localhost:8080/swagger/index.html
```

### Docker

```bash
# Запустить всё через Docker
make docker-up
# или
docker compose up -d

# Посмотреть логи
make docker-logs

# Остановить
make docker-down
```

## 🧪 Тестирование

```bash
# Запустить тесты
make test

# Тесты с отчётом покрытия
make test-coverage
```

## 📊 Качество кода

```bash
# Запустить линтер
make lint

# Установить инструменты разработки
make tools
```

## 🚀 Деплой на Render.com

1. Форкните репозиторий
2. Создайте аккаунт на [render.com](https://render.com)
3. Создайте новый Web Service из вашего репозитория
4. Render автоматически определит Dockerfile
5. Добавьте PostgreSQL database
6. Настройте переменные окружения:
   - `DATABASE_URL` (автоматически из базы данных)
   - `JWT_SECRET` (сгенерируйте надёжный ключ)
   - `ENVIRONMENT=production`

Или используйте Blueprint (`render.yaml`) для автоматического деплоя.

## 📁 Структура проекта

```
.
├── cmd/
│   └── api/
│       └── main.go          # Точка входа
├── internal/
│   ├── config/              # Конфигурация
│   ├── domain/              # Модели и ошибки
│   ├── handler/             # HTTP handlers
│   │   └── middleware/      # Middleware (auth, cors)
│   ├── repository/          # Data access layer
│   │   └── postgres/        # PostgreSQL implementation
│   └── service/             # Business logic
├── docs/                    # Swagger документация
├── .github/
│   └── workflows/           # GitHub Actions
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

## 👥 Команда

Проект разработан в рамках курса "Проектирование человеко-машинных интерфейсов" БГУ.

## 📄 Лицензия

MIT License
