
def api_demo(id, i,send_back):
    # check input here    
    if "Data" not in i:
        return send_back(id, {"Err":"missing parameter Data", "Vector": None})
    # your logic here
    send_back(id,{"Result": "demo ok"})
