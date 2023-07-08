### server, python example:
```
class service_textToMp3(Service):
    def __init__(self):
        Service.__init__(self,"redis.vm:6379/0")
    def check_input_item(self, i):
            if "BackTo" not in i:
                return False
            return True
    def process(self,items):
        #your logic here
        for i in items:
            self.send_back(i,{"Result":input.value})
service_textToMp3().start()
```
