# Сервис проведения тендеров

# Использование

**Сборка**

Из корневой директории проекта:

`$ docker build -t avito_tenders .`

**Запуск**

Из корневой директории проекта:

`$ docker run -d -p 8080:8080 --name avito_container avito_tenders`

**Сборка и запуск с БД**

Можно запустить для тестирования проект с помощью Dockerfile_DB, находящегося также в корне проекта. Он запустит проект вместе 
с базой данных Postgres в докере, выполнит скрипт db/dbex.sql с тестовыми данными и таблицами.

Из корневой директории проекта:

`$ docker build -f Dockerfile_DB -t avito_tenders_db .`

`$ docker run -d -p 8080:8080 --name avito_container_db avito_tenders_db`

**Переменные окружения**

Все переменные окружения по умолчанию заданы в докерфайлах. Можно их изменить либо внутри файлов, либо при запуске
докер контейнера, пример:

`$ docker run -d -p 8080:8080 -e POSTGRES_DATABASE=tenders -e POSTGRES_PORT=5433 --name avito_container_db avito_tenders_db`

**Запуск линтера**

Из корневой директории проекта:

`$ golangci-lint run -c .golangci.yml`

