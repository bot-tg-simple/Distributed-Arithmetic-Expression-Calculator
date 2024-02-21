Запуск системы

Можно перенести содержимое такой командой в терминале (VsCode)-> git clone https://github.com/bot-tg-simple/Distributed-Arithmetic-Expression-Calculator
Затем переходим в директорию Distributed-Arithmetic-Expression-Calculator командой -> cd Distributed-Arithmetic-Expression-Calculator
А затем вызваем следующее -> go run main.go (чтобы подгрузить библиотеки перед запуском -> go mod download; ещё нужна библиотека github.com/Knetic/govaluate(она установится с командой go mod download, но если что-то не так -> go get github.com/Knetic/govaluate)
После открываем файл index.html(можно открыть папку и щёлкнуть пару раз на index.html) и вычисляем согласно дальнейшей инструкции. По всем вопросам -> контакты в ТГ ниже(в самом конце).

Установите необходимые зависимости:(если не получится через go mod download)

go get github.com/Knetic/govaluate
go get -u github.com/gin-gonic/gin
go get -u github.com/jmoiron/sqlx
go get -u github.com/mattn/go-sqlite3

Все библиотеки, которые нужны:

"github.com/Knetic/govaluate"
	_ "github.com/Knetic/govaluate"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"


Запустите приложение:

go run main.go(если, конечно, вы назвали файл программы main.go)

Система будет доступна по адресу: http://localhost:8080

Использование HTTP API

Добавление выражения

curl -X POST http://localhost:8080/expression -d "expression=2+2"

Получение списка выражений

curl http://localhost:8080/expressions

Получение информации о конкретном выражении

curl http://localhost:8080/expression/{id}

Получение списка операций

curl http://localhost:8080/operations

Получение списка ресурсов

curl http://localhost:8080/resources

Получение текущих выполняющихся операций

curl http://localhost:8080/task

Обновление результата выражения

curl -X POST http://localhost:8080/result -d "id={id}" -d "resultStr=4.0"

Обновление длительности операции по имени

curl -X POST http://localhost:8080/operation-duration -d "operation=+" -d "duration=15"

Проверка подключения к серверу

Система автоматически проверяет подключение к серверу. Если подключение потеряно, будет возвращен статус ошибки.

Примеры запросов и ответов

Добавление выражения: /expression

1) curl -X POST http://localhost:8080/expression -d "expression=2%2B(2*2)/(2-1)"

2) curl -X POST http://localhost:8080/expression -d "expression=2%2B2" (то есть, 2+2 => 2%2B2 в URL-кодировке("+" = %2B))

3) curl -X POST http://localhost:8080/expression -d "expression=2-2"

4) curl -X POST http://localhost:8080/expression -d "expression=2/2"

5) curl -X POST http://localhost:8080/expression -d "expression=2*3"

Пример ответа:

1) {"id":"b330606e-3f2b-4ed4-b301-f79dd0abafcd"}

2) {"id":"9f45db54-adea-4b53-bd1d-52c35fb06f80"}

3) {"id":"6dbc2c24-559c-4120-8522-5d2129f51032"}

4) {"id":"60660f75-4b7b-4656-94a8-fec942b13e57"}

5) {"id":"098748f4-8ea5-4a58-b236-75ac82173256"}

Поиск по ID выражения: /expression/:id

1) curl -X GET http://localhost:8080/expression/b330606e-3f2b-4ed4-b301-f79dd0abafcd

2) curl -X GET http://localhost:8080/expression/9f45db54-adea-4b53-bd1d-52c35fb06f80

3) curl -X GET http://localhost:8080/expression/6dbc2c24-559c-4120-8522-5d2129f51032

4) curl -X GET http://localhost:8080/expression/60660f75-4b7b-4656-94a8-fec942b13e57

5) curl -X GET http://localhost:8080/expression/098748f4-8ea5-4a58-b236-75ac82173256

Пример ответа:

1) {"ID":"b330606e-3f2b-4ed4-b301-f79dd0abafcd","Expression":"2+(2*2)/(2-1)","Status":"completed","CreatedAt":"2024-02-18T19:51:07.007244+03:00","UpdatedAt":"2024-02-18T19:51:07.007244+03:00","Result":6}

2) {"ID":"9f45db54-adea-4b53-bd1d-52c35fb06f80","Expression":"2+2","Status":"completed","CreatedAt":"2024-02-18T19:52:16.637351+03:00","UpdatedAt":"2024-02-18T19:52:16.637351+03:00","Result":4}

3) {"ID":"6dbc2c24-559c-4120-8522-5d2129f51032","Expression":"2-2","Status":"completed","CreatedAt":"2024-02-18T19:55:23.13525+03:00","UpdatedAt":"2024-02-18T19:55:23.13525+03:00","Result":0}

4) {"ID":"60660f75-4b7b-4656-94a8-fec942b13e57","Expression":"2/2","Status":"completed","CreatedAt":"2024-02-18T19:56:49.731848+03:00","UpdatedAt":"2024-02-18T19:56:49.731849+03:00","Result":1}

5) {"ID":"098748f4-8ea5-4a58-b236-75ac82173256","Expression":"2*3","Status":"completed","CreatedAt":"2024-02-18T20:00:00.151546+03:00","UpdatedAt":"2024-02-18T20:00:00.151547+03:00","Result":6}

Запрос на /operations:

Пример запроса:

curl -X GET http://localhost:8080/operations

Пример ответа:

[{"Name":"+","Duration":5,"StartTime":"0001-01-01T00:00:00Z","Status":""},{"Name":"-","Duration":10,"StartTime":"0001-01-01T00:00:00Z","Status":""},{"Name":"*","Duration":15,"StartTime":"0001-01-01T00:00:00Z","Status":""},{"Name":"/","Duration":20,"StartTime":"0001-01-01T00:00:00Z","Status":""}]


Запрос на /resources:

Пример запроса:

curl -X GET http://localhost:8080/resources

Пример ответа:

[{"Name":"Resource1","Operation":"+","Duration":0},{"Name":"Resource2","Operation":"-","Duration":0},{"Name":"Resource3","Operation":"*","Duration":0},{"Name":"Resource4","Operation":"/","Duration":0}]

Запрос на /expressions(получение всех выражений):

Пример запроса:

curl -X GET http://localhost:8080/expressions

Пример ответа:

[{"ID":"b330606e-3f2b-4ed4-b301-f79dd0abafcd","Expression":"2+(2*2)/(2-1)","Status":"completed","CreatedAt":"2024-02-18T19:51:07.007244+03:00","UpdatedAt":"2024-02-18T19:51:07.007244+03:00","Result":6},{"ID":"9f45db54-adea-4b53-bd1d-52c35fb06f80","Expression":"2+2","Status":"completed","CreatedAt":"2024-02-18T19:52:16.637351+03:00","UpdatedAt":"2024-02-18T19:52:16.637351+03:00","Result":4},{"ID":"6dbc2c24-559c-4120-8522-5d2129f51032","Expression":"2-2","Status":"completed","CreatedAt":"2024-02-18T19:55:23.13525+03:00","UpdatedAt":"2024-02-18T19:55:23.13525+03:00","Result":0},{"ID":"60660f75-4b7b-4656-94a8-fec942b13e57","Expression":"2/2","Status":"completed","CreatedAt":"2024-02-18T19:56:49.731848+03:00","UpdatedAt":"2024-02-18T19:56:49.731849+03:00","Result":1},{"ID":"e4a0ba89-7fae-4008-a4dc-76543e68d8ad","Expression":"2*2","Status":"completed","CreatedAt":"2024-02-18T19:58:29.936426+03:00","UpdatedAt":"2024-02-18T19:58:29.936426+03:00","Result":4},{"ID":"098748f4-8ea5-4a58-b236-75ac82173256","Expression":"2*3","Status":"completed","CreatedAt":"2024-02-18T20:00:00.151546+03:00","UpdatedAt":"2024-02-18T20:00:00.151547+03:00","Result":6}]

Запрос на /task(получение выражений, которые выполняются):

Пример запроса:

curl -X GET http://localhost:8080/task

Пример ответа:

{"ID":"40d428e2-d983-46c5-a849-57098b278777","Expression":"2-2","Status":"pending","Duration":1000}


Запрос на /result:

Пример запроса: 

curl -X POST \
  http://localhost:8080/result \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'id=5c4d91fc-dab8-4b5a-a419-b4a024edce31' \
  -d 'resultStr=3.14'

Пример ответа:

{"message":"Result updated successfully"}

После делаем на /expression/:id(для проверки):

Запрос:

curl -X GET http://localhost:8080/expression/5c4d91fc-dab8-4b5a-a419-b4a024edce31

Ответ:
{"ID":"5c4d91fc-dab8-4b5a-a419-b4a024edce31","Expression":"2*3","Status":"completed","CreatedAt":"2024-02-18T20:16:52.084201+03:00","UpdatedAt":"2024-02-18T20:16:52.084201+03:00","Result":3.14}


Запрос на /operation-duration:

Пример запроса:

curl -X POST \
  http://localhost:8080/operation-duration \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'operation=%2B' \
  -d 'duration=10'

В URL-кодировке("+" = %2B)

Пример ответа:

{"message":"Operation duration updated successfully"}

После делаем на /operations(для проверки):

Запрос:

curl -X GET http://localhost:8080/operations

Ответ:

[{"Name":"+","Duration":10,"StartTime":"0001-01-01T00:00:00Z","Status":""},{"Name":"-","Duration":10,"StartTime":"0001-01-01T00:00:00Z","Status":""},{"Name":"*","Duration":15,"StartTime":"0001-01-01T00:00:00Z","Status":""},{"Name":"/","Duration":20,"StartTime":"0001-01-01T00:00:00Z","Status":""}]

Также сервер мониторит время последнего ping'а, и если прошло больше минуты с последнего, то станет доступен запрос на /ping(также в ТЕРМИНАЛ будет вывод "Connection to the server is lost: the item will stop displaying after 1 minute"):

Запрос:

curl -X GET http://localhost:8080/ping

Ответ:

{"error":"Connection to the server is lost: the item will stop displaying after 1 minute"}

Но после получения этого ответа сервер отключится(типо пропало с ним соединение спустя минуту)
P.S. его можно будет снова запустить.


Также есть интерфейс(файлы - index.html, index1.html, index2.html, index3.html, styles.css, script.js, script1.js, script2.js) В браузере надо открывать файл index.html (например, в папке проекта, куда вы перенесёте содержимое, кликнуть на index.html и он откроется в браузере)

На странице Калькулятор надо вводить выражение в поле и затем нажать Enter:
лучше вводить простые выражения типа: 2-2; 2*2
деление(/) и сложение лучше делать через curl запрос и не забыть про плюс писать так в URL-кодировке "+" = "%2B"

! На ноль не стоит делить на странице и делать curl запросы(но если сделали curl запрос, то результат независимо от операции можно посмотреть на странице Поиск по ID введя ID в поле и нажав Enter)

Ниже выводится ID выражения, его можно скопировать и вставить на странице Поиск по ID в поле. Также нажать Enter и ниже будет выводиться вся информация о выражении.

Страница Настройка расчёта пока не работает.

На странице Вычислительные ресурсы можно нажать на кнопку Check и увидеть есть ли соединение с сервером. Если не прошло минуты с последнего пинга то будет зеленый статус, а если прошло больше минуты будет уведомление об остановке сервера и красный статус, но и сервер, конечно, завершит работу.

!!!

Схема в файле SCHEME.md

!!!

Если что-то пошло не так, то лучше остановить систему в терминале (VSCode -> ^c (то есть control + C(Mak); Ctrl + C(Windows)));
затем удалить базу данных (expressions.db) и перезапустить систему (go run main.go).

Мои контакты в ТГ на случай вопросов -> https://t.me/ZOV22_2023


