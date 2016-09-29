#! /usr/bin/env python

import argparse
import random
import socket
import sys
import time

def http_request(sock, i):
    msg = "GET index.html HTTP/1.0"
    try:
        sock.sendall(msg)
    except Exception as e:
        print("%d: Send error: %s" % (i, e))
        sys.exit(1)

    received = 0
    while True:
        data = sock.recv(4096)
        received += len(data)
        if len(data) == 0:
            break
    print("Received %d bytes" % received)

parser = argparse.ArgumentParser()
parser.add_argument('-c', '--conn', type=int, default=256,
                    help="Number of concurrent connections")
parser.add_argument('-t', '--tries', type=int, default=10,
                    help="Number of tries to connect")
parser.add_argument('-p', '--port', type=int, default=8080,
                    help="Port to connect to")
parser.add_argument('-s', '--server', default="localhost",
                    help="server to connect to")
args = parser.parse_args()

sockets = []

def connect():
    addr = (args.server, args.port)
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        sock.connect(addr)
        return sock
    except Exception as e:
        sys.stderr.write("%d: %s\n" % (i, e))
        sys.exit(1)

print("Creating %d connections" % args.conn)
for i in range(args.conn):
    sockets.append(connect())

# Issue a HTTP request on a random socket
for i in range(args.tries):
    idx = random.randrange(0, len(sockets))
    http_request(sockets[idx], idx)
    sockets[idx].close()
    sockets[idx] = connect()

for sock in sockets:
    sock.close()
