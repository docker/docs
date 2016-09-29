#! /usr/bin/env python
# Simple UDP echo server

import os
import socket
import sys
import time

host = sys.argv[1]
port = int(sys.argv[2])

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
server_address = (host, port)
sock.bind(server_address)

while True:
    print ('Listening for data on %s:%d' % (host, port))
    data, address = sock.recvfrom(65535)

    print ('received %s bytes from %s' % (len(data), address))
    print (data)

    if data:
        sent = sock.sendto(data, address)
        print ('sent %s bytes back to %s' % (sent, address))
