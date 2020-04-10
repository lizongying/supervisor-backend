# supervisor-backtend
本项目可以和
[supervisor-frontend](https://github.com/lizongying/supervisor-frontend)
结合使用

## dev 
```
go run supervisor.go -c ./conf/dev.yml
```

## prod
```
go build
./supervisor -c ./conf/prod.yml
```
![preview](./Screenshot.png)

