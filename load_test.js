import http from 'k6/http';

export let options = {
    vus: 500, // Количество виртуальных пользователей
    duration: '30s', // Длительность теста
    rps: 1000 // Ограничение на 1000 запросов в секунду
};

export default function () {
    let params = {
        headers: {
            'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk3NTQxMDQsInVzZXJuYW1lIjoi0JjQstCw0L0ifQ.2GJ33MTOUX2jfA2b82kbClRUaOnmlK7iGr6Tef1yktM'
        }
    };
    http.get('http://localhost:8080/api/info', params);
}