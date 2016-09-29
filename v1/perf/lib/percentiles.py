#!/usr/bin/env python

import sys

ping_file = sys.argv[1]
with open(ping_file) as f:
    content = [ int(line) for line in f.readlines() ]

content.sort()
length = len(content)

if len(sys.argv) > 2:
    only = int(sys.argv[2])
    print content[int(only * length / 100.)]
else:
    for percentile in range(10,100,10):
        print "%d\t%d" % (percentile, content[int(percentile * length / 100.)])
