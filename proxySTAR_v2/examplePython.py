import datetime

"""d = datetime.datetime.utcnow() # <-- get time in UTC
e = datetime.datetime.utcfromtimestamp(1)
f = datetime.datetime.utcnow() + datetime.timedelta(days=4)
g = datetime.datetime.utcnow() + datetime.timedelta(days=0)
print d
print e
print f
print g
time = f.isoformat('T') + 'Z'
time2 = g.isoformat('T') + 'Z'
#time = time + 
print time
print time2
"""
def exeTime(lengthInDays=0):
    d = datetime.datetime.utcnow() + datetime.timedelta(days=lengthInDays)
    return d.isoformat('T') + 'Z'

a = exeTime()
print (type(a) is str)
print exeTime(365)

