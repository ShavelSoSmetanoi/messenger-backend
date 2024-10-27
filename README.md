# Название проекта

Описание проекта, его предназначение и основные функции.

## Установка и настройка

### 1. Клонирование репозитория

Сначала клонируйте репозиторий проекта и перейдите в его директорию:

```bash
git clone https://github.com/ShavelSoSmetanoi/messenger-backend.git
cd messenger-backend
```

### 2. Установка зависимостей

Проект написан на Go и использует модули для управления зависимостями. Чтобы установить все зависимости, выполните:

```bash
go mod tidy
```

### 3. Поднятие Docker Compose
Проект использует Docker Compose для работы с необходимыми сервисами, такими как база данных. Чтобы поднять все сервисы, выполните:

```bash
docker-compose up -d
```

### 4. Установка Goose для управления миграциями
Для управления миграциями используется утилита Goose. Установите её, выполнив следующую команду:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```
Убедитесь, что $GOPATH/bin добавлен в переменную окружения PATH, чтобы команду goose можно было использовать из любого места.

Для применения всех миграций к базе данных выполните команду:

```bash
goose -dir migrations postgres "user=ваш_пользователь password=ваш_пароль dbname=имя_вашей_базы sslmode=disable" up
```
По дефоту будет:
```bash
goose -dir migrations postgres "user=myuser password=mypassword dbname=mydatabase sslmode=disable" up
```