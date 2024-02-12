import msgpack
from .rds import rds
#__all__ = ["Do","DoAt"]


def DoAt(ServiceKey, paramIn, timeAt):
    global rds
    _fields = {"data": msgpack.packb(paramIn)}
    if timeAt != 0:
        _fields.update({"timeAt": str(timeAt)})

    # ServiceKey
    if ServiceKey[:4] != "api:":
        ServiceKey = "api:" + ServiceKey
    cmd_id = rds.xadd(name=ServiceKey, fields=_fields, id="*", maxlen=4096)
    if timeAt != 0 or cmd_id == None:
        return None
    # BLPop 返回结果 [key1,value1,key2,value2]
    # cmd.Val() is the stream id, the result will be poped from the list with this id
    results = rds.blpop(keys=[cmd_id], timeout=20)
    if results==None or len(results) <2:
        return None
    return msgpack.unpackb(results[1])


def Do(ServiceKey, paramIn):
    return DoAt(ServiceKey, paramIn, 0)
