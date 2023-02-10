import redis
rds = redis.Redis(connection_pool=redis.ConnectionPool(
    host="keydb2.vm", port=6379, db=0, decode_responses=False))