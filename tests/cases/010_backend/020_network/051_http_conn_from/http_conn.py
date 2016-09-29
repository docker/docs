#! /usr/bin/env python

import argparse
import random
import os
import subprocess
import sys
import threading
import time

status = 0
failed = []

# spread the love
urls = [
    "http://www.google.com/",
    "http://www.docker.com/",
    "https://www.facebook.com/",
    "http://www.amazon.com/",
    "https://azure.microsoft.com/",
    "https://www.yahoo.com/index.html",
    "http://www.bbc.co.uk/",
    "https://www.bing.com/",
    "https://www.youtube.com/",
    "https://www.baidu.com/",
    "https://www.twitter.com/",
    "https://www.live.com/",
    "https://www.linkedin.com/",
    "https://www.ebay.com/",
    "https://www.netflix.com/",
    "https://apple.com",
    "https://PayPal.com",
    "https://msn.com",
]

devnull = open(os.devnull, 'wb')


class thread(threading.Thread):
    def __init__(self, tid, iterations):
        threading.Thread.__init__(self)
        self.tid = tid
        self.iters = iterations

    def run(self):
        global results
        for i in range(self.iters):
            idx = random.randrange(0, len(urls))
            ret = subprocess.call(["curl", "-4", "-s", urls[idx]],
                                  stdout=devnull)
            if not ret == 0:
                print("ERR[%02d:%d] %d %s" % (self.tid, i, ret, urls[idx]))
                failed[self.tid] += 1

            time.sleep(0)

parser = argparse.ArgumentParser()
parser.add_argument('-c', '--conn', type=int, default=10,
                    help="Number of connections per thread")
parser.add_argument('-t', '--threads', type=int, default=10,
                    help="Number of threads")
args = parser.parse_args()

# create threads
threads = []
for i in range(args.threads):
    failed.append(0)
    threads.append(thread(i, args.conn))

# run threads
for t in threads:
    t.start()

# wait for them to finish
for t in threads:
    t.join()

devnull.close()

tot_fail = sum(failed)
print("Connections attempts: %d  Failed: %d" %
      ((args.conn * args.threads), sum(failed)))

if ((tot_fail * 100.0) / (args.conn * args.threads)) > 10.0:
    print("More than 10% of connections failed.")
    sys.exit(1)
