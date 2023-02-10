from inspect import getmembers, isfunction,ismodule,getmodulename
from api import api,Do,DoAt

#register all api functions
import api_funcs
for name, apiFunction in getmembers(api_funcs,isfunction):
    if  apiFunction.__name__.startswith("api_"):
         api.new_service(apiFunction)

import threading
#new thread to report states
threading.Thread(target=api.reportStates).start()
#new thread to receive jobs
threading.Thread(target=api.receiveJobs).start()


import time
time.sleep(5)
#test Do demo
print("result of test demo api:",Do("demo",{"Data":"我的国家"}))

#sleep max 68 years
time.sleep(2**32)
