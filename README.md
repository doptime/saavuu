## saavuu, the most concise, redis based web-server framework
    the name saavuu borrow from "杀悟",means kill bad wisdom。 I hate bad tools.
### major advantages on API design
* APIs you defined, support both monolithic and microservice architecture. perfect!
* very simple to define and use API. just see the golang example demo.
* No API version related problem. Just see web client demo.
* Very easy to upgrade API , just change the data structure. no extra schema definition needed.
* You don't need to write any CREATE GET PUT or DELETE  Logic. Just use redis to query，modify or delete. That means most CURD can be done at frontend, needs no backend job.
* You can focus on operations with multiple data logic only.  We call it "API".
    saavuu will put API data to redis stream, and the API receive and process the stream data.
* redis pipeline  brings high batch process performance.  
### major advantages on Data Op
* Using most welcomed redis compatible db. no database but redis compatible KEYDB. With flash storage supportion, KEYDB brings both memory speed and whole disk capacity
* Very Easy to define and access data. see keyInDemo.HSET(req.Id, req) in golang example.
 - Schema data is adopted to keep maintain-bility. Easy to upgrade data structure.
### other features
* Use msgpack to support structure data by default. Easily to upgrade data sturecture.
* All HTTP requests are transferd as binary msgpack data. It's compact and fast.
* allow specify Content-Type in web client.
* allow specify response fields in web client to reduce web traffic
* support JWT for authorization
* fully access control
* support CORS
  

## demo usage
### server, golang example:
```
package main

import (
	"github.com/yangkequn/saavuu/api"
)

type InDemo struct {
	Data   []uint16
	Id   string `msgpack:"alias:JWT_id"`
}
var keyInDemo = data.NewStruct[*InDemo]()
//define api with input/output data structure
ApiDemo,_=api.Api(func(req *InDemo) (ret string, err error) {
    // your logic here
    if req.Id == "" || len(req.Data) == 0 {
        return nil, saavuu.ErrInvalidInput
    }
    keyInDemo.HSET(req.Id, req)
    return `{data:"ok"}`, nil
})

// calling api
func main() {
    //your logic here
    ApiDemo.Call(&InDemo{Data:[]uint16{1,2,3},Id:"1234567890"})
}
```

### web client, javascript /typescript example:
```
HGET("UserInfo", id).then((data) => {
    //your logic here
})
```

## about configuration 
    saavuu reads configuration from enviroment variables. Make sure enviroment variables are added to your IDE (launch.json for vs code) or docker. 
    these are the default example:
```
    "RedisAddress_PARAM": "127.0.0.1:6379",
    "RedisPassword_PARAM": "",
    "RedisDb_PARAM": "0",
    "RedisAddress_DATA": "127.0.0.1:6379",
    "RedisPassword_DATA": "",
    "RedisDb_DATA": "0",
    "JWTSecret": "WyBJujUQzWg4YiQqLe9N36DA/7QqZcOkg2o=",
    "JWT_IGNORE_FIELDS": "iat,exp,nbf,iss,aud,sub,typ,azp,nonce,auth_time,acr,amr,at_hash,c_hash,updated_at,nonce,auth_time,acr,amr,at_hash,c_hash,updated_at",
    "CORS": "*",
    "MaxBufferSize": "3145728",
    "AutoPermission": "true",
```