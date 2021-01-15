# supervisor-backtend

![preview](./Screenshot.png)
base on gin

## 【推荐】和supervisor-frontend结合使用，

[supervisor-frontend](https://github.com/lizongying/supervisor-frontend)

### dev

```
go run supervisor.go -c ./dev.yml
```

### prod

linux build

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o supervisor
```

mac build

```
go build -o supervisor
```

run

```
./supervisor -c example.yml
```

