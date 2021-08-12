# gin-api-server-template

## Prerequisites:
- Install swaggo [link](https://github.com/swaggo/swag)

## Generate swagger docs:
Use following code to generate swagger documents.
```
> swag init -g .\internal\api\api.go --output .\api\docs\
```

## Run gin server:
```
> go run .\cmd\api-server\main.go -cfgpath .\configs\config.ini
```