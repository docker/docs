#!/usr/bin/env python

import time, os, sys

start = time.time ()
os.system(" ".join(sys.argv[2:]))
total = time.time () - start

f = open(sys.argv[1], "w")
f.write(str(total))
f.close()
