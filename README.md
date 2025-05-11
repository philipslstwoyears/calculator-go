

# gRPC Calculator

## Описание

Этот проект представляет собой API для калькулятора с возможностью регистрации и входа пользователей. API позволяет пользователю выполнять вычисления через HTTP-ручки. Все внутренние сервисы общаются друг с другом через gRPC, а вычисления выполняются асинхронно.

## Особенности

* **Асинхронные вычисления**: Вычисления выполняются асинхронно, чтобы не блокировать основное приложение.
* **gRPC**: Все внутренние сервисы, включая обработку вычислений, используют gRPC для эффективной коммуникации.
* **Регистрация и аутентификация**: Поддерживается регистрация пользователей и аутентификация по логину и паролю.

---

## 🚀 Запуск сервера

1. **Склонируйте репозиторий**
   ```sh
   git clone https://github.com/philipslstwoyears/calculator-go.git
   cd calculator-go
   ```
   
2. **Установите зависимости:**
   ```sh
   go mod tidy
   ```

3. **Запустите сервер**
   ```sh
    go run ./cmd/main.go
   ```
   Сервер будет доступен на `http://localhost:8080`.

---

## Ручки API

### 1. Регистрация пользователя

**Метод**: `POST`

**URL**: `/api/v1/register`

**Описание**: Регистрирует нового пользователя с логином и паролем.

#### Пример запроса:

```cmd
curl -X POST http://localhost:8080/api/v1/register -H "Content-Type: application/json" -d "{\"login\": \"testuser\", \"password\": \"12345\"}"
```

#### Ответ (успешный запрос):

```json
{
  "id": 1
}
```

#### Ошибка (пользователь уже зарегистрирован):

```cmd
curl -X POST http://localhost:8080/api/v1/register -H "Content-Type: application/json" -d "{\"login\": \"testuser\", \"password\": \"12345\"}"
```

Ответ:

```json
{
  "error": "user is already registered"
}
```

---

### 2. Вход в систему

**Метод**: `POST`

**URL**: `/api/v1/login`

**Описание**: Выполняет вход в систему, возвращая cookie с идентификатором пользователя.

#### Пример запроса:

```cmd
curl -X POST http://localhost:8080/api/v1/login -H "Content-Type: application/json" -d "{\"login\": \"testuser\", \"password\": \"12345\"}"
```

#### Ответ (успешный запрос):

```json
{
  "id": 1
}
```

#### Ошибка (неправильный логин или пароль):

```cmd
curl -X POST http://localhost:8080/api/v1/login -H "Content-Type: application/json" -d "{\"login\": \"wronguser\", \"password\": \"wrongpass\"}"
```

Ответ:

```json
{
  "error": "wrong login or password"
}
```

---

### 3. Выполнение вычисления

**Метод**: `POST`

**URL**: `/api/v1/calculate`

**Описание**: Выполняет вычисление математического выражения асинхронно.

#### Пример запроса:

```cmd
curl -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -H "Cookie: id=1" -d "{\"expression\": \"2+2*2\"}"
```

#### Ответ (успешный запрос):

```json
{
  "id": 123
}
```

#### Ошибка (отсутствует cookie с id):

```cmd
curl -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d "{\"expression\": \"2+2*2\"}"
```

Ответ:

```json
{
  "error": "http: named cookie not present"
}
```

---

### 4. Получение списка выражений

**Метод**: `GET`

**URL**: `/api/v1/expressions`

**Описание**: Получает все выражения для текущего пользователя.

#### Пример запроса:

```cmd
curl -X GET http://localhost:8080/api/v1/expressions -H "Cookie: id=1"
```

#### Ответ (успешный запрос):

```json
[
  {
    "id": 1,
    "expression": "2+2*2",
    "status": "completed",
    "result": 6
  },
  {
    "id": 2,
    "expression": "3+5*5",
    "status": "completed",
    "result": 28
  }
]
```

#### Ошибка (отсутствует cookie с id):

```cmd
curl -X GET http://localhost:8080/api/v1/expressions
```

Ответ:

```json
{
  "error": "http: named cookie not present"
}
```

---

### 5. Получение конкретного выражения

**Метод**: `GET`

**URL**: `/api/v1/expressions/{id}`

**Описание**: Получает результат вычисления для конкретного выражения по его ID.

#### Пример запроса:

```cmd
curl -X GET http://localhost:8080/api/v1/expressions/1 -H "Cookie: id=1"
```

#### Ответ (успешный запрос):

```json
{
  "id": 1,
  "expression": "2+2*2",
  "status": "completed",
  "result": 6
}
```

#### Ошибка (выражение не найдено):

```cmd
curl -X GET http://localhost:8080/api/v1/expressions/100 -H "Cookie: id=1"
```

Ответ:

```json
{
  "error": "expression not found"
}
```

---

## Важная информация

* **Асинхронная обработка**: Все вычисления в калькуляторе выполняются асинхронно. Когда вы отправляете запрос на выполнение вычисления, запрос сразу возвращает ID задачи. Результат можно получить позже, используя ID этого вычисления.

* **gRPC**: Внутренние сервисы для вычислений и управления пользователями взаимодействуют через gRPC, что позволяет минимизировать время отклика и повысить производительность.

---

## Заключение

Это приложение предоставляет простой и мощный API для регистрации пользователей, выполнения вычислений и управления историей вычислений. Все запросы выполняются эффективно с использованием современных технологий, таких как gRPC и асинхронная обработка.

---
> Автор: *philipslstwoyears*  
> 📌 *[GitHub](https://github.com/philipslstwoyears)*