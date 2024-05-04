# URL Shortener

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)

Микросервис для генерации коротких ссылок.
Сгенерированные ссылки закрепляются за пользователем, в ответе в хедере проставляется Cookie

Что умеет:
- сохранять одиночные/множество ссылок
- выдача оригинальных ссылок
- просмотр всех сгенерированных ссылок пользователя (нужен Cookie)
- удаление сгенерированных ссылок пользователем (нужен Cookie)
- поддерживает 3 варианта хранения (postgres, файловое хранение, in-memory)

Используемые технологии:
- PostgreSQL (как одно из хранилищ)
- Docker (для запуска сервиса)
- Chi (роутинг, middlewares)
- golang-migrate/migrate (миграции БД)
- pgx (драйвер для работы с Postgres)
- jwt (генерация токенов в Cookie)
- zerolog (логирование)
- testify (тесты)

# Getting Started

Для запуска сервиса необходимо создать и заполнить `.env` файл, а также иметь Docker и docker-compose на машине

# Usage

запуск сервиса
```sh
make up
```

остановка сервиса
```sh
make down
```


## Examples

Примеры запросов:
- [Генерация короткой ссылки через `Content-Type: text/plain`](#gen_short_simple)
- [Переход по короткой ссылке](#get_original_url)
- [Генерация короткой ссылки через `Content-Type: application/json`](#gen_short_json)
- [Массовая генерация коротких ссылок](#gen_short_batch)
- [Просмотр всех созданных пользователем ссылок](#get_user_urls)
- [Удаление коротких ссылок пользователем](#delete_user_urls)

### Генерация короткой ссылки (простая) <a name="gen_short_simple"></a>

Запрос:
```curl
curl -v -d "<your_url>" http://localhost:8080
```
Ответом будет сгенерированная короткая ссылка

### Переход по короткой ссылке <a name="get_original_url"></a>

Запрос:
```curl
curl -v http://localhost:8080/<your_short_url>
```
Ответом будет редирект на оригинальную ссылку

### Генерация короткой ссылки (json) <a name="gen_short_json"></a>

Запрос:
```curl
curl -v -H "Content-Type: application/json" \
  -d '{"url": "<your_url>"}' \
http://localhost:8080/api/shorten
```
Ответом будет сгенерированная короткая ссылка

### Массовая(Batch) генерация коротких ссылок <a name="gen_short_batch"></a>

Запрос:
```curl
curl -v -H "Content-Type: application/json" \
  -d '[{"original_url": "<your_url>", "correlation_id": "<uuid>"}]'
http://localhost:8080/api/shorten/batch
```
Ответом будет массив коротких ссылок с correlation_id

### Получение всех ссылок пользователем <a name="get_user_urls"></a>

Для доступа к ссылкам необходимо подставить Cookie выданный ранее

Запрос:
```curl
curl -v -H "Cookie:<your_cookie>" http://localhost:8080/api/user/urls
```
Ответом будет все созданные ссылки+короткие для пользователя

### Удаление ссылок пользователем <a name="delete_user_urls"></a>

Для доступа необходимо подставить Cookie выданный ранее

Запрос:
```curl
curl -v -X DELETE \
  -H "Cookie:<your_cookie>" \
  -H "Content-Type: application/json" \
  -d '["short_url", "short_url"]' \
http://localhost:8080/api/user/urls
```
