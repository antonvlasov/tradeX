# Сведения о программе
- Язык реализации go
- Использование базы данных sqlite
- Запуск через docker
# Инструкция для запуска
Необходимо иметь: docker; git.
 - Скачать содежримое репозитория командой git clone https://github.com/antonvlasov/tradeX
 - Перейти в скачанную директорию : cd tradeX
 - Построить контейнер командой docker build -t tradex-docker .
 - Запустить контейнер: docker run --rm --name tradex-instance -p 1778:1778 tradex-docker
# Инструкция по использованию
 С запущенным контейнером можно взаимодействовать через порт 1778 с помощью RPC 1.0 запросов: 
- Добавление события: {"method":"Database.AddEvent","params":[{"date":"YYYY-MM-DD","views":vies,"clicks":count,"cost":cost}],"id":id}
Ответ вида: {"id":id,"result":0,"error":null}
- Просмотр событий за период с сортировкой результата:
{"method":"Database.SelectStats","params":[{"from":"YYYY-MM-DD","to":"YYYY-MM-DD","sortby":"sortby"}],"id":0}
 Сортировка возможна по следующим полям: date,views,clicks,cost.
Ответ вида: {"id":id,"result":[{"date":"YYYY-MM-DD1","views":views1,"clicks":clicks1,"cost":cost1,"cpc":cpc1,"cpm":cpm1},{"date":"YYYY-MM-DD2","views":views2,"clicks":clicks2,"cost":cost2,"cpc":cpc2,"cpm":cpm2}],"error":null}
Пример ошибки при неправильном формате даты:
{"id":0,"result":null,"error":"invalid dates"}
- Сброс статистики: {"method":"Database.Clear","params":[0],"id":id}, здесь значение "params":[0] должно быть таким всегда.
Ответ: {"id":0,"result":0,"error":null}
# Примеры запросов
- Через консоль:
запрос:
echo -e '{"method":"Database.AddEvent","params":[{"date":"2021-02-03","views":27000,"clicks":900240,"cost":3000}],"id":0}' | nc localhost 1778
ответ:
{"id":0,"result":0,"error":null}
запрос: echo -e '{"method":"Database.SelectStats","params":[{"from":"2021-02-03","to":"2021-02-04","sortby":"date"}],"id":0}' | nc localhost 1778
ответ:
{"id":0,"result":[{"date":"2021-02-03","views":27000,"clicks":900240,"cost":3000,"cpc":0.0033324446814182885,"cpm":111.1111111111111}],"error":null}
- С использованием go:
Создание клиента:
``` var client *rpc.Client
	var err error
	for client, err = jsonrpc.Dial("tcp", "localhost:1778"); err != nil; client, err = jsonrpc.Dial("tcp", "localhost:1778") {
		fmt.Println("connecting...")
	}
```
Запросы:
```
    var mock *int
	err = client.Call("Database.Clear", 0, &mock) // Запрос на очистку статистики
	err = client.Call("Database.AddEvent", events[i], mock) // Запрос на добавление статистики
	err = client.Call("Database.SelectStats", database.StatRequest{date1, date2, sortby}, &stats) // Запрос на выбор статистики
```
# Тесты

Покрытие тестами 80%, cover-файл coverage.out