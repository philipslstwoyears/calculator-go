### 🧮 Async Calculator API

**Async Calculator API** – это асинхронный калькулятор, который принимает математические выражения, обрабатывает их в фоне и предоставляет результат по уникальному идентификатору запроса.

---

## 🚀 Запуск сервера

1. **Склонируйте репозиторий**
   ```sh
   git clone https://github.com/your-repo/calculator-go.git
   cd calculator-go
   ```

2. **Запустите сервер**
   ```sh
   go run main.go
   ```
   Сервер будет доступен на `http://localhost:8080`.

---

## 📌 API Эндпоинты

### ➕ 1. Вычисление выражения
**`POST /api/v1/calculate`**

🔹 **Описание:** Добавляет выражение в очередь на вычисление.  
🔹 **Тело запроса (JSON):**
   ```json
   {
   "expression": "5+5"
}
   ```
🔹 **Ответ (JSON):**
   ```json
   {
   "id": 1
}
   ```
📌 **Важно:** Результат вычисления можно получить по ID через `/api/v1/expressions/{id}`.

✅ **Сценарий "Успешное выполнение"**
- Запрос:
  ```sh
  curl -X POST "http://localhost:8080/api/v1/calculate" \
       -H "Content-Type: application/json" \
       -d '{"expression": "5+5"}'
  ```
- Ответ:
  ```json
  {
      "id": 1
  }
  ```

❌ **Сценарий "Ошибка в выражении"**
- Запрос:
  ```sh
  curl -X POST "http://localhost:8080/api/v1/calculate" \
       -H "Content-Type: application/json" \
       -d '{"expression": "5/0"}'
  ```
- Ответ:
  ```json
  {
      "error": "Деление на ноль"
  }
  ```

❌ **Сценарий "Неверный формат запроса"**
- Запрос:
  ```sh
  curl -X POST "http://localhost:8080/api/v1/calculate" \
       -H "Content-Type: application/json" \
       -d '{"wrong_field": "5+5"}'
  ```
- Ответ:
  ```json
  {
      "error": "Некорректный запрос"
  }
  ```

---

### 📜 2. Получение списка всех выражений
**`GET /api/v1/expressions`**

🔹 **Описание:** Возвращает все выражения, которые были отправлены на вычисление.  
🔹 **Пример ответа:**
   ```json
   [
       {
           "id": 1,
           "expression": "5+5",
           "status": "Ok",
           "result": 10
       },
       {
           "id": 2,
           "expression": "10/2",
           "status": "В обработке"
       }
   ]
   ```

✅ **Сценарий "Успешное выполнение"**
- Запрос:
  ```sh
  curl -X GET "http://localhost:8080/api/v1/expressions"
  ```
- Ответ:
  ```json
  [
      {
          "id": 1,
          "expression": "5+5",
          "status": "Ok",
          "result": 10
      },
      {
          "id": 2,
          "expression": "10/2",
          "status": "В обработке"
      }
  ]
  ```

❌ **Сценарий "Список пуст"**
- Ответ:
  ```json
  []
  ```

---

### 🔍 3. Получение результата по ID
**`GET /api/v1/expressions/{id}`**

🔹 **Описание:** Возвращает информацию о конкретном вычислении.  
🔹 **Пример ответа для ID = 1:**
   ```json
   {
       "id": 1,
       "expression": "5+5",
       "status": "Ok",
       "result": 10
   }
   ```
📌 Если выражение еще в обработке, `result` будет отсутствовать.

✅ **Сценарий "Выражение вычислено"**
- Запрос:
  ```sh
  curl -X GET "http://localhost:8080/api/v1/expressions/1"
  ```
- Ответ:
  ```json
  {
      "id": 1,
      "expression": "5+5",
      "status": "Ok",
      "result": 10
  }
  ```

✅ **Сценарий "Выражение еще вычисляется"**
- Запрос:
  ```sh
  curl -X GET "http://localhost:8080/api/v1/expressions/2"
  ```
- Ответ:
  ```json
  {
      "id": 2,
      "expression": "10/2",
      "status": "В обработке"
  }
  ```

❌ **Сценарий "Выражение не найдено"**
- Запрос:
  ```sh
  curl -X GET "http://localhost:8080/api/v1/expressions/999"
  ```
- Ответ:
  ```json
  {
      "error": "Выражение с таким ID не найдено"
  }
  ```

---

## 🛠 Описание технологий

- **Язык:** Go
- **Маршрутизация:** `gorilla/mux`
- **Логирование и обработка ошибок:** Middleware
- **Асинхронная обработка:** Каналы (`chan`), горутины

---

## 📌 Пример запроса с `cURL`
### 🔹 Отправка выражения на вычисление:
```sh
curl -X POST "http://localhost:8080/api/v1/calculate" \
     -H "Content-Type: application/json" \
     -d '{"expression": "5+5"}'
```
### 🔹 Получение списка всех выражений:
```sh
curl -X GET "http://localhost:8080/api/v1/expressions"
```
### 🔹 Получение результата вычисления по ID:
```sh
curl -X GET "http://localhost:8080/api/v1/expressions/1"
```

---

## 🏁 Заключение
Этот проект демонстрирует **асинхронную обработку вычислений** в Go, используя **goroutines, каналы и middleware**. 🚀

> Автор: *philipslstwoyears*  
> 📌 *[GitHub](https://github.com/philipslstwoyears)*