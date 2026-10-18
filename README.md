##Сервис сегментирования пользователей
HTTP-сервис для управления пользователями и сегментами.

Позволяет создавать/удалять пользователей, управлять сегментами (CRUD),
добавлять пользователей в сегменты и удалять из них, получать список сегментов
конкретного пользователя, а также случайным образом распределять сегмент
на заданный процент всех пользователей (например, «30% пользователей в MAIL_GPT»)

##Стек 
- **Go 1.26** + Gin 
- **PostgreSQL** 
- **Docker** + docker-compose
- **golang-migrate** 


## Возможности

### Пользователи
- Создание пользователя
- Получение списка всех пользователей
- Получение пользователя по ID
- Удаление пользователя по ID 

### Сегменты
- Создание сегмента с проверкой уникальности имени
- Получение списка всех сегментов
- Получение сегмента по ID
- Изменение имени сегмента 
- Удаление сегмента по ID 

### Связи пользователь <-> сегмент
- Добавление пользователя в сегмент (с проверкой существования обоих)
- Удаление пользователя из сегмента
- Получение списка сегментов конкретного пользователя

### Распределение
- Случайное добавление сегмента заданному проценту пользователей
  (например, 30% — с перемешиванием и защитой от дублей)


1. git clone https://github.com/slinkeres/segment-service
2. cd segment-service
3. Скопируй env.example в .env
4. Запускаем контейнеры
```bash
docker-compose up --build
```
## Примеры запросов

### Пользователи

**Создать пользователя:**

```bash
curl -X POST http://localhost:8080/users
```

Ответ:
```json
{
  "message": "user created",
  "id": 1
}
```

**Список пользователей:**

```bash
curl http://localhost:8080/users
```

Ответ:
```json
{
  "result": [
    { "id": 1, "created_at": "2026-10-01T12:00:00Z" },
    { "id": 2, "created_at": "2026-10-01T12:01:00Z" }
  ]
}
```

**Пользователь по ID:**

```bash
curl http://localhost:8080/users/1
```

**Удалить пользователя:**

```bash
curl -X DELETE http://localhost:8080/users/1
```

### Сегменты

**Создать сегмент:**

```bash
curl -X POST http://localhost:8080/segments \
  -H "Content-Type: application/json" \
  -d '{"name":"MAIL_GPT"}'
```

Ответ:
```json
{
  "message": "segment created",
  "id": 1,
  "name": "MAIL_GPT"
}
```

При дубликате имени — `409 Conflict`:
```json
{ "message": "segment name exists" }
```

**Список сегментов:**

```bash
curl http://localhost:8080/segments
```

**Изменить имя сегмента:**

```bash
curl -X PATCH http://localhost:8080/segments/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"MAIL_GPT_V2"}'
```

**Удалить сегмент:**

```bash
curl -X DELETE http://localhost:8080/segments/1
```

### Связи

**Добавить пользователя в сегмент:**

```bash
curl -X POST http://localhost:8080/user-segments/1/users/2
```

Ответ:
```json
{ "message": "user add to segment" }
```

При дубликате — `409 Conflict`:
```json
{ "message": "user already exists in segment" }
```

**Удалить пользователя из сегмента:**

```bash
curl -X DELETE http://localhost:8080/user-segments/1/users/2
```

**Сегменты пользователя:**

```bash
curl http://localhost:8080/users/2/segments
```

Ответ:
```json
{
  "result": [
    { "id": 1, "name": "MAIL_GPT", "created_at": "2026-10-01T12:00:00Z" }
  ]
}
```

### Распределение

**Добавить сегмент `id=1` для 30% всех пользователей:**

```bash
curl -X POST http://localhost:8080/segments/1/apply-percent \
  -H "Content-Type: application/json" \
  -d '{"num":30}'
```

Поле `num` — процент от 1 до 100. Ответ:

```json
{
  "message": "users was added to segment",
  "result": [
    { "id": 3, "created_at": "2026-10-01T12:00:00Z" },
    { "id": 7, "created_at": "2026-10-01T12:00:00Z" }
  ]
}
```
