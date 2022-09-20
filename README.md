# saavuu
most concise, convinient, redis based microservice framework

# specification
* specify content-type in header,if response type is not json, then return raw data
* use JWT for authorization, JWT field "id" will replace @me in key or field
* when get request, if field is not exist, then return all the hash key list for given key

# congifuration is read from environment variables
### make sure enviroment variables are added to your IDE or docker. 
### for example, if you are using vscode, add this to your launch.json
```
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${fileDirname}/main.go",
            "env": {
                "REDIS_ADDR_PARAM": "docker.vm:6379",
                "REDIS_PASSWORD_PARAM": "",
                "REDIS_DB_PARAM": "0",
                "REDIS_ADDR_DATA": "docker.vm:6379",
                "REDIS_PASSWORD_DATA": "",
                "REDIS_DB_DATA": "0",
                "JWT_SECRET": "WyBJujUQzWg4YiQqLe9N36DA/7QqZcOkg2o=",
                "CORS": "*",
                "SAAVUU_CONFIG_KEY": "saavuu_service_config",
                "MAX_BUFFER_SIZE": "3145728",
            },
        }
```