# Songs Library

Проект **Songs Library** - это приложение на Go, которое предоставляет функциональность для работы с библиотекой песен. Приложение использует базу данных PostgreSQL для хранения данных о песнях.

## Требования

Для запуска проекта на вашем локальном компьютере вам потребуется:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)

## Установка и запуск

Следуйте этим шагам, чтобы запустить проект локально:

1. **Клонируйте репозиторий**:

   ```bash
   git clone https://github.com/Russo2642/songs_lib.git
   cd songs_lib 
   
2. **Создайте .env файл**: \
   Создайте файл .env в корневой директории проекта и добавьте следующие переменные окружения:
    - DB_HOST=db
    - DB_PORT=5432
    - DB_NAME=postgres
    - DB_USERNAME=postgres
    - DB_PASSWORD=postgres

    - API_URL=апи_для_детали_песни
    
    - POSTGRES_DB=postgres
    - POSTGRES_USER=postgres
    - POSTGRES_PASSWORD=postgres \
    
3. **Запустите приложение и базу данных**:

   ```bash
   docker-compose up --build

4. **Откройте приложение**: \
   Приложение доступно по адресу: http://localhost:8000. \
   Для просмотра swagger документации откройте - http://localhost:8000/swagger/index.html#/
