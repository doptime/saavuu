# saavuu, the most concise, redis based framework
* All HTTP requests are transferd as binary msgpack data. It's compact and fast.
* No API version related problem. Just use redis api.
* Use msgpack to support structure data by default. Easily to upgrade data sturecture.
* Use no database but KEYDB which is redis compatible. Flash storage supportion brings both memory speed and whole disk capacity
* You don't need to write any GET Logic. Just use redis to query.
* You can focus only and alway on Modification logic (POST PUT DEL) logic. 
* You can use any programming language you like. python or go or not.
* redis pipeline  brings heavy batch process performance.  

# demo usage
```
package main

import (
	"github.com/yangkequn/saavuu"
	"github.com/yangkequn/saavuu/rCtx"
)

type Input struct {
	data   []uint16
	JWT_id string
}

func init() {
	saavuu.NewService("demo", 128, func(dc *rCtx.DataCtx, pc *rCtx.ParamCtx, parmIn map[string]interface{}) (data map[string]interface{}, err error) {
		var req *Input = &Input{}
		if err = saavuu.MapsToStructure(parmIn, req); err != nil {
			return nil, err
		} else if req.JWT_id == "" || len(req.data) == 0 {
			return nil, saavuu.ErrInvalidInput
		}
		// your logic here
		data = map[string]interface{}{"data": "ok"}
		return data, nil
	})
}
```

# abstract    
    for query, saavuu will query  result from redis and return to client.
    for modification, saavuu will put request data to redis queue, and the service listening the queue will process the data.
    saavuu means kill bad wisdom, which borrow from "杀悟"。 I hate bad tools.


# feature
* specify content-type in header,response fields etc. in client side
* support JWT for authorization
* convient privacy control for redis batch operation.


# about configuration 
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