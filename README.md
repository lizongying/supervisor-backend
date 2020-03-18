# supervisor-backtend
本项目可以和
[supervisor-frontend](https://github.com/lizongying/supervisor-frontend)
结合使用

## dev 
```
export GIN_MODE=debug 
go run supervisor.go -c /Users/lizongying/IdeaProjects/supervisor-backend/conf/dev.yml
```

## prod
```
export GIN_MODE=release
./supervisor -c /Users/lizongying/IdeaProjects/supervisor-backend/conf/prod.yml
```

