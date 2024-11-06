# messenger-backend

## Установка и настройка

### 1. Клонирование репозитория

Сначала клонируйте репозиторий проекта и перейдите в его директорию:

```bash
git clone https://github.com/ShavelSoSmetanoi/messenger-backend.git
cd messenger-backend
```

### 2. Создание файлов .env
Для корректной работы приложения вам нужно создать файл .env как в корне проекта, так и в папке deployments, содержащий необходимые настройки для работы с базой данных и JWT.

Шаг 1: В корне проекта создайте файл .env со следующим содержимым:

```bash
JWT_SECRET=your_jwt_secret_here
DB_USER=myuser
DB_PASSWORD=mypassword
DB_NAME=mydatabase
DB_HOST=my_postgres_container 
DB_PORT=5432
```
Замените your_jwt_secret_here на ваш секретный ключ для JWT, а также параметры подключения к базе данных (если они отличаются от значений по умолчанию).

Шаг 2: В папке deployments создайте файл .env:

```bash
PG_USER=myuser
PG_PASSWORD=mypassword
PG_NAME=mydatabase
```

### 3. Установка зависимостей

Проект написан на Go и использует модули для управления зависимостями. Чтобы установить все зависимости, выполните:

```bash
go mod tidy
```

### 4. Поднятие Docker Compose
Проект использует Docker Compose для работы с необходимыми сервисами, такими как база данных. Чтобы поднять все сервисы, выполните:

```bash
docker-compose -f deployments/docker-compose.yml up -d
```

### 5. Установка Goose для управления миграциями
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