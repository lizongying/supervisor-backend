### dev 
```
export GIN_MODE=debug 
go run supervisor.go -c /Users/lizongying/IdeaProjects/supervisor/conf/dev.yml
```

### prod
```
export GIN_MODE=release
./supervisor -c /Users/lizongying/IdeaProjects/supervisor/conf/prod.yml
```

