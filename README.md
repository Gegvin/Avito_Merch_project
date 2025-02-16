
Avito Merch Project

Это мой проект для стажировки. 

Я написал его на Go, PostgreSQL , JWT  и Docker Compose. Проект позволяет сотрудникам приобретать мерч за монетки и  переводить их друг другу. 

(Интерфес самый простой html код. Просто изходя из задачи, сотрудникм было неудобно бы пользоваться командной строкой, поэтому сделал + веб)


Что я реализовал

API эндпоинты:

/api/auth – Авторизация 

curl -X POST http://localhost:8080/api/auth \
     -H "Content-Type: application/json" \
     -d '{"username": "your_username", "password": "your_password"}'

/api/register – регистрация нового пользователя

curl -X POST http://localhost:8080/api/register \
     -H "Content-Type: application/json" \
     -d '{"username": "new_username", "password": "new_password"}'

     В ответах вам сообщет ваши токены
     
/api/info – получение информации о кошельке, инвентаре и истории транзакций

curl -X GET http://localhost:8080/api/info \
     -H "Authorization: Bearer <JWT_TOKEN>"
     
     
/api/sendCoin – перевод монет другому сотруднику

curl -X POST http://localhost:8080/api/sendCoin \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDQwMTMzMDIsInVzZXJuYW1lIjoieW91cl91c2VybmFtZSJ9.someSignature" \
     -d '{"toUser": "recipient_username", "amount": 50}'
     
/api/buy/{item} – покупка мерча по уникальному названию

curl -X GET http://localhost:8080/api/buy/t-shirt \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDQwMTMzMDIsInVzZXJuYW1lIjoieW91cl91c2VybmFtZSJ9.someSignature"
     
Веб-интерфейс:

Страницы для входа, регистрации, просмотра кошелька, покупки мерча и перевода монет.
	
 Тестирование:
 
 •	Я написал юнит-тесты для сервисов и репозитория, а также интеграционные тесты для проверки работы API. В итоге общее покрытие тестами получилось больше 40%.

Какие проблемы пришлось решать

	•	Подключение к базе данных:
 
Сначала я столкнулся с проблемами подключения, так как база в Docker пробрасывалась на порт 5433, а в моём коде использовался стандартный 5432. (как я понял мой стандартный постгрес сидит на 5432 и у них конфликт) Я настроил переменную окружения DB_CONFIG так, чтобы использовать порт 5433 – теперь всё работает как надо.

	•	JWT и защита эндпоинтов:
 
Защищённые маршруты возвращали 401, потому что JWT не попадал туда (функция getUserDataFromToken искала его в cookie). Мне пришлось изменить интеграционные тесты,  чтобы передавать токен через cookie, теперь запросы проходят корректно.

Почему-то иногда докер не запускает контейнер с самим сервером (помогает перезапуск)


Как запустить проект через Docker Compose

Убедитесь, что Docker установлен

переходите в директорию проекта 

docker-compose up -d

Приложение будет доступно по адресу http://localhost:8080.




Как запустить тесты
	•	Для юнит-тестов и покрытия:

go test ./... -v -coverprofile=coverage.out

для интеграционных

DB_CONFIG="host=localhost port=5433 user=postgres password=postgres dbname=avito_merch sslmode=disable" JWT_SECRET="mysecret" go test ./test/integration -v


Для нагрузочного тестирования я использовал к6 и скрипт (Скрипт на js писал не сам)



     execution: local
        script: load_test.js
        output: -

     scenarios: (100.00%) 1 scenario, 500 max VUs, 1m0s max duration (incl. graceful stop):
              * default: 500 looping VUs for 30s (gracefulStop: 30s)




     data_received..................: 5.6 MB 184 kB/s
     data_sent......................: 7.5 MB 247 kB/s
     http_req_blocked...............: avg=10.77µs  min=0s     med=3µs      max=22.62ms  p(90)=6µs      p(95)=9µs     
     http_req_connecting............: avg=4.56µs   min=0s     med=0s       max=17.71ms  p(90)=0s       p(95)=0s      
     http_req_duration..............: avg=6.96ms   min=493µs  med=789µs    max=492.68ms p(90)=990µs    p(95)=5.9ms   
       { expected_response:true }...: avg=6.88ms   min=493µs  med=789µs    max=492.68ms p(90)=983µs    p(95)=5.55ms  
     http_req_failed................: 0.10%  31 out of 30499
     http_req_receiving.............: avg=27.47µs  min=4µs    med=22µs     max=17.02ms  p(90)=35µs     p(95)=47µs    
     http_req_sending...............: avg=12.1µs   min=1µs    med=9µs      max=8.04ms   p(90)=16µs     p(95)=21µs    
     http_req_tls_handshaking.......: avg=0s       min=0s     med=0s       max=0s       p(90)=0s       p(95)=0s      
     http_req_waiting...............: avg=6.92ms   min=464µs  med=757µs    max=492.65ms p(90)=925.2µs  p(95)=5.8ms   
     http_reqs......................: 30499  1000.025051/s
     iteration_duration.............: avg=495.86ms min=37.2ms med=499.98ms max=972.89ms p(90)=500.17ms p(95)=501.15ms
     iterations.....................: 30499  1000.025051/s
     vus............................: 500    min=500         max=500
     vus_max........................: 500    min=500         max=500


running (0m30.5s), 000/500 VUs, 30499 complete and 0 interrupted iterations
default ✓ [======================================] 500 VUs  30s




