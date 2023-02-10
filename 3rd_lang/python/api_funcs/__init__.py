#from .api_demo import api_demo
#from .api_wav2vec2 import api_wav2vec

#rather than import all api_*.py files one by one, we auto import do it
#scan all api_*.py files in api_funcs folder,and import them
import os
for file in os.listdir(os.path.dirname(__file__)):
    if file.startswith("api_") and file.endswith(".py"):
        import_name = file[:-3]
        import_module = "." + import_name
        import_statement = "from %s import %s" % (import_module, import_name)        
        exec(import_statement)
