# saavuu
most concise, convinient, redis based microservice framework

# specification
* specify content-type in header,if response type is not json, then return raw data
* use JWT for authorization, JWT field "id" will replace @me in key or field
* when get request, if field is not exist, then return all the hash key list for given key