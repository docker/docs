#! /usr/bin/env python
# Simple file transfer utility (nc in python)

import os
import socket
import sys
import time

host = sys.argv[1]
port = int(sys.argv[2])
filename = sys.argv[3]

fsz = os.stat(filename).st_size
f = open(filename, 'rb')

tries = 0
while True:
    try:
        s = socket.socket()        
        s.connect((host, port))
        break
    except Exception as e:
        tries += 1
        if tries >= 10:
            sys.stderr.write("Can't connect after %d tries\n" % tries)
            sys.exit(1)
        time.sleep(1)
        
# connected
sz = 0
b = f.read(2048)
while b:
    s.sendall(b)
    sz += len(b)
    b = f.read(2048)
f.close()
s.close()

print("Filesize: %d transferred: %d" % (fsz, sz))
if not fsz == sz:
    sys.exit(1)
else:
    sys.exit(0)
