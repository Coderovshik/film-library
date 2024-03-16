# :books: Фильмотека 

### Запуск приложения

Может занять некоторое время из-зи прогонки тестов и хелфчека базы данных.

```bash
make compose-build-up
```

### Документация

Документация расположена по пути [`api/openapi.yml`](https://github.com/Coderovshik/film-library/blob/master/api/openapi.yml). Для просмотра можно использовать [Swagger Editor](https://editor.swagger.io/).

### Тесты

Запуск тестов в docker-контейнере осуществляется командой:

```bash
make test-docker
```

### Подробности реализации

- **Стек:** Go, PostgreSQL, Docker, OpenAPI 3.0
- **Авторизация:** JWT
- **Поиск фильмов:** Сортировка и фильтрация осуществляется с пмомщью query-параметров
- **Администратор:** Логин - admin_user, пароль - admin_password
- **Данные:** По умолчанию в базу данных загружен небольшой объем mock-данных (за подробностями обращайтесь к [файлам миграций](https://github.com/Coderovshik/film-library/tree/master/internal/db/migrations))