echo -e '{"method":"Database.AddEvent","params":[{"date":"2021-02-04","views":27000,"clicks":9000,"cost":3000}],"id":0}' | nc localhost 1778
{"id":0,"result":0,"error":null}


echo -e '{"method":"Database.SelectStats","params":[{"from":"2021-02-03","to":"2021-02-04","sortby":"date"}],"id":0}' | nc localhost 1778
{"id":0,"result":[{"date":"2021-02-04","views":27000,"clicks":9000,"cost":3000,"cpc":0.3333333333333333,"cpm":111.1111111111111}],"error":null}

(base) [antonvlasov@PC32 ~]$ echo -e '{"method":"Database.Clear","params":[0],"id":0}' | nc localhost 12345
{"id":0,"result":0,"error":null}
(base) [antonvlasov@PC32 ~]$ echo -e '{"method":SelectStats","params":[{"from":"2021-02-03","to":"2021-02-04","sortby":"date"by":"date"}],"id":0}' | nc localhost 1778
{"id":0,"result":[],"error":null}
