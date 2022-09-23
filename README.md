# saavuu, the most concise, convinient, redis based microservice framework
* 1. You dont't need to write any HTTP API, api address is fixed. you will never need to consider api version related problem.
* 2. You need no database but redis. Use api to easily save and load any data structure. Using redis compitable db that support flash storage brings both performance persistance and capacity.
* 3. You will never need to write any GET Logic. Saavuu query query result from redis directly. It's alway possible and cheap.
* 4. You can focus only and alway on Modification logic (POST PUT DEL) logic. 
* 5. Saavuu allow Plug and Play services very easily without restart any other service. Very suitable for team development.
* 6. All HTTP requests are transferd as binary msgpack data. It's compact and fast. And decode to data structure in your logic directly.
* 7. redis pipeline 's batch read enabling every high volume request to be processed simultaneously.  

# abstract    
    saavuu take data from client and throw to redis queue, and the service listening the queue will process the data. all api from client or backend work in this way.
    saavuu means kill bad wisdom, which borrow from "杀悟"。 I hate bad tools.

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
                "JWT_IGNORE_FIELDS": "iat,exp,nbf,iss,aud,sub,typ,azp,nonce,auth_time,acr,amr,at_hash,c_hash,updated_at,nonce,auth_time,acr,amr,at_hash,c_hash,updated_at",
                "CORS": "*",
                "SAAVUU_CONFIG_KEY": "saavuu_service_config",
                "MAX_BUFFER_SIZE": "3145728",
            },
        }
```
# JWT is read from "Authorization" fields in header. all jwt fields are and sent to service, except those defined in JWT_IGNORE_FIELDS
    if you do not want to 
# Never use COOKIE, always use JWT

# privacy control for GET method.
    Only Key start with Upper case will be returned, Other wise access is not allowed
    This idea is borrow from golang, where only public method start with Upper case will be exported
    
    if filed in get contains @field, then the "field" will be replaced with the value of "field" in JWT

# Annotation for api.tsx
    QueryFields should be "*" or "field1,field2,field3" 