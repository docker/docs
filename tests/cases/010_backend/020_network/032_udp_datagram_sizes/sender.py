#! /usr/bin/env python
# Simple UDP ping/pong program

import os
import socket
import sys
import time
import threading
import Queue

host = sys.argv[1]
port = int(sys.argv[2])
size = sys.argv[3]

s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

message = "X" * int(size)

success = Queue.Queue()

def receiver():
    count = 0
    while True:
        try:
            print "Listening for a response..."
            data, server = s.recvfrom(int(size))
            if int(size) == len(data):
                print "Response has the correct size (%s)" % size
                count = count + 1
                if count >= 2:
                    success.put(True)
                    break
            else:
                print "Expected length %s, received length %d" % (size, len(data))
        except Exception as e:
            print "Ignoring %s\n" % e
            time.sleep(1)

t = None

tries = 0
while tries < 30:
    if not success.empty():
        print "Success: message received"
        sys.exit(0)
    print "Sending %s byte UDP message to %s:%s" % (size, host, port)
    s.sendto(message, (host, port))
    if t is None:
        t = threading.Thread(target=receiver)
        t.daemon = True
        t.start()
    time.sleep(1)
    tries = tries + 1
print "No response from server"
sys.exit(1)
