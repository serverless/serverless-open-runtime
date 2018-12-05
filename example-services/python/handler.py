import json

def hello(event):
  return {'body': json.dumps(event)}
