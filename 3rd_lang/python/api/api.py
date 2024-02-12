# provide service via redis API
from .rds import rds
import msgpack
import datetime
import time
#__all__ = ['api']


class api():
    service_names = []
    batch_size = 64
    ApiFunc = {}
    task_num_in_60s = {}

    def new_service(fuc):
        service_name = fuc.__name__.split("_")[1]
        if not service_name.startswith("api:"):
            service_name = "api:" + service_name
        api.service_names.append(service_name)
        api.task_num_in_60s[service_name] = 0
        api.ApiFunc[service_name] = fuc
        print(f"service {service_name } added")

    def reportStates():
        now = datetime.datetime.now()
        while True:
            time.sleep(60)
            now = datetime.datetime.now()
            for service_name in api.service_names:
                print(
                    f"{now} py service {service_name} job rcved in 60s:{api.task_num_in_60s[service_name]}")
                api.task_num_in_60s[service_name] = 0

    def XGroupCreate():
        global rds
        for serviceName in api.service_names:
            # create group0 for each stream
            try:
                cmd = rds.xgroup_create(
                    name=serviceName, groupname="group0", id="$", mkstream=True)
                if cmd == True:
                    print("create groupd: "+cmd)
            except Exception as e:
                print(str(e))
            # create consumer saavuu for each stream
            try:
                cmd = rds.xgroup_createconsumer(
                    name=serviceName, groupname="group0", consumername="saavuu")
                if cmd == True:
                    print("create consumer: "+cmd)
            except Exception as e:
                print(str(e))

    def receiveJobs():
        global rds
        api.XGroupCreate()
        streams = dict(zip(api.service_names, [">"]*len(api.service_names)))
        while True:
            try:
                # read group using group0, python as consumer, '>' as the last id, 20s timeout, 64 batch size
                ret = rds.xreadgroup(groupname="group0", consumername="saavuu",
                                    streams=streams, count=api.batch_size, block=2000, noack=True)
                if ret == None or len(ret) == 0:
                    continue
                print("cmd: "+str(ret))
                for stream in ret:
                    stream_name = stream[0]
                    apiName = stream_name.decode("utf-8")
                    for message in stream[1]:
                        id = message[0]
                        _messege = message[1]
                        param = msgpack.unpackb(_messege[b'data'])
                        if b"timeAt" in _messege:
                            timeAtStr = _messege[b"timeAt"]
                            api.delayTaskAddOne(
                                apiName, timeAtStr, param)
                        else:
                            api.ApiFunc[apiName](id, param, api.send_back)
                        api.task_num_in_60s[apiName]+=1
            except Exception as e:
                print(str(e))

    def send_back(id, output, use_single_float=False, use_bin_type=True):
        global rds
        if id == None:
            return
        pipe = rds.pipeline()
        # https://msgpack-python.readthedocs.io/en/latest/api.html
        packer = msgpack.Packer(
            use_single_float=use_single_float, use_bin_type=use_bin_type)
        pipe.rpush(id, packer.pack(output))
        pipe.expire(id, 6)
        pipe.execute()

    def delayTaskAddOne(serviceName, timeAtStr, bytesValue):
        global rds
        rds.hset(serviceName+":delay", timeAtStr, bytesValue)
        api.delayTaskDoOne(serviceName, timeAtStr)

    def delay_task_do_one(serviceName, timeAtStr):
        global rds
        nowUnixMilliSecond = time.time() * 1000
        timeAtUnixMilliSecond = int(timeAtStr)
        time.sleep((timeAtUnixMilliSecond-nowUnixMilliSecond) / 1000)
        rds.hget(serviceName+":delay", timeAtStr)
        api.rds.hdel(serviceName+":delay", timeAtStr)
        if bytes != None:
            api.ApiFun[serviceName](None, bytes, api.send_back)

    def delay_tasks_load():
        services = api.ApiFun.keys()
        for service in services:
            timeAtStrs = api.rds.hkeys(service+":delay")
            for timeAtStr in timeAtStrs:
                api.delay_task_do_one(service, timeAtStr)
