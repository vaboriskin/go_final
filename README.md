# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

реализованы все задачи со звездочкой, кроме установки потоврений по неделям и месяцам.


.env_example
-------------------
TODO_PASSWORD=12345
TODO_DBFILE=./scheduler.db
TODO_PORT=7540
-------------------

для получения токена используем команду
curl -X POST http://localhost:7540/api/signin \
     -H "Content-Type: application/json" \
     -d '{"password":"12345"}'

    
сборка докер 
docker build -t go_final .

можно создать папку например data куда будет сохраняться файл для БД
запуск контейнера в таком случае 
docker run --rm --name go_final `
  -p 7540:7540 `
  -e TODO_PORT=7540 `
  -e TODO_DBFILE=/data/scheduler.db `
  -e TODO_PASSWORD=12345 `
  -v "${PWD}\data:/data" `
  go_final