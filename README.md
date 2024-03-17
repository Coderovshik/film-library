# :clapper: Фильмотека 

### Запуск приложения

Может занять некоторое время из-за прогонки тестов и хелфчека базы данных.

```bash
make compose-build-up
```

### Документация

Документация предоставляется в двух форматах: HTML и YAML. HTML-документ доступен по адресу [**localhost:8080/docs/html**](http://localhost:8080/docs/html) (**Обратите внимание**: из-за того, что html-генератор не в полной мере поддерживает спецификацию OpenAPI 3.0 - не все элементы отображаются корректно - за полной документацией обращайтесь к yaml-файлу). YAML-документ доступен по адресу [**localhost:8080/docs/yaml**](http://localhost:8080/docs/yaml).

### Тесты

Запуск тестов в docker-контейнере осуществляется командой:

```bash
make test-docker # прогон тестов и покрытие в docker-контейнере
```

Если вам доступна среда исполнения Go, можете воспользоваться командами:

```bash
make test # прогон тестов
make test-coverage # покрытие
```

### Подробности реализации

- **Стек:** Go, PostgreSQL, Docker, OpenAPI 3.0
- **Авторизация:** JWT
- **Поиск фильмов:** Сортировка и фильтрация осуществляется с помощью query-параметров
- **Администратор:** Логин - admin_user, пароль - admin_password
- **Данные:** По умолчанию в базу данных загружен небольшой объем mock-данных (За подробностями обращайтесь к [файлам миграций](https://github.com/Coderovshik/film-library/tree/master/internal/db/migrations))
- **Сервер:** По умолчанию сервер доступен по адресу [**localhost:8080**](http://localhost:8080) (Порт можно поменять в `config/.env`).