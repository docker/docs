#! /usr/bin/env python

import argparse
import ssl
import sys
import threading
import time
if sys.version_info < (3, 0):
    import urllib2 as urlreq
else:
    import urllib.request as urlreq

status = 0
failed = []

class thread(threading.Thread):
    def __init__(self, tid, iterations, url):
        threading.Thread.__init__(self)
        self.tid = tid
        self.iters = iterations
        self.req = urlreq.Request(url)

    def run(self):
        global results
        for i in range(self.iters):
            try:
                r = urlreq.urlopen(self.req)
                page = r.read()
            except Exception as e:
                print("%02d:%d: %s" % (self.tid, i, e))
                # XXX Hopefully this does not require locking :)
                failed[self.tid] += 1
            # yield
            time.sleep(0)

# Very skanky hack to disabled SSL Certificate authentication
ssl._create_default_https_context = ssl._create_unverified_context

parser = argparse.ArgumentParser()
parser.add_argument('-c', '--conn', type=int, default=1000,
                    help="Number of connections")
parser.add_argument('url', nargs=1,
                    help="URL to fetch")
args = parser.parse_args()

# pick number of thread based on number of connections
num_thds = 20 if args.conn > 1000 else 10


iterations = args.conn / num_thds

# create threads
threads = []
for i in range(num_thds):
    failed.append(0)
    threads.append(thread(i, iterations, args.url[0]))

# run threads
for t in threads:
    t.start()

# wait for them to finish
for t in threads:
    t.join()

tot_fail = sum(failed)
print("Connections attempts: %d  Failed: %d" % (args.conn, sum(failed)))

if ((tot_fail * 100.0) / args.conn) > 10.0:
    print("More than 10% of connections failed.")
    sys.exit(1)
