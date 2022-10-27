# hezzl_test5

Чтобы запустить, надо

1) Запустить Postgres и накатить миграции

    - goose postgres "user=postgres password=postgres dbname=postgres host=localhost" up

2) Запустить Redis

    - systemctl start redis

3) Запустить Nats

    - запуск из докер файла "docker run -p 4222:4222 -ti nats:latest"

4) Запустить CLickHouse и создать таблицы(накатить миграции)

    - clickhouse-client 

