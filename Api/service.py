# provide service via redis API
import redis
import msgpack
import datetime
import time


class Service():
    def __init__(self, service_name, batch_size, host, port=6379, db=0):
        connpool = redis.ConnectionPool(
            host=host, port=port, db=db, decode_responses=False)
        self.rds = redis.Redis(connection_pool=connpool)
        self.batch_size = batch_size
        self.service_name = service_name
        if not service_name.startswith("svc:"):
            self.service_name = "svc:" + service_name
        self.task_num_in_60s = datetime.datetime.now().minute << 32
        print(f"service {self.service_name } initialed, batch_size:{self.batch_size} host:{host} port:{port} db:{db}")

    def get_task(self):
        while True:
            # https://stackoverflow.com/questions/20621775/pop-multiple-values-from-redis-data-structure-atomically
            pipe = self.rds.pipeline()
            pipe.lrange(self.service_name, start=0, end=self.batch_size-1)
            pipe.ltrim(self.service_name, start=self.batch_size, end=-1)
            data = pipe.execute()
            if data == None or len(data[0]) == 0:
                time.sleep(0.05)
                continue

            self.task_num_in_60s += len(data[0])
            now = datetime.datetime.now()
            if (self.task_num_in_60s >> 32) != now.minute:
                print(
                    f"{now} py service {self.service_name} job rcved in 60s:{self.task_num_in_60s&0xFFFFFFFF}")
                self.task_num_in_60s = now.minute << 32
            ds = [msgpack.unpackb(k) for k in data[0]]

            # property check, all properties should provided
            ds1 = [i for i in ds if self.check_input_item(i)]
            if len(ds1) > 0:
                return ds1
            print(
                f"warning! service {self.service_name} ,input data corrupt :{ds}")

    def send_back(self, i, output, use_single_float=False, use_bin_type=True):
        pipe = self.rds.pipeline()
        # https://msgpack-python.readthedocs.io/en/latest/api.html
        packer = msgpack.Packer(
            use_single_float=use_single_float, use_bin_type=use_bin_type)
        pipe.rpush(i["BackTo"], packer.pack(output))
        pipe.expire(i["BackTo"], 6)
        pipe.execute()

    def start(self):
        print(f"service {self.service_name} started")
        while True:
            inputs = self.get_task()
            self.process(inputs)
## use like this:
# class Service_xxx(Service):
#     def __init__(self):
#         Service.__init__(self,"xxx",1,"docker.vm",6379,15,)

#     def process(self,items):
#         # your logic here
#         for i in items:
#             self.send_back(i,{"Result":input.value})
#     def check_input_item(self, i):
#             if "BackTo" not in i:
#                 return False
#             return True
